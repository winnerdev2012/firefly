// Copyright © 2021 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package data

import (
	"context"
	"fmt"
	"testing"

	"github.com/kaleido-io/firefly/internal/config"
	"github.com/kaleido-io/firefly/mocks/databasemocks"
	"github.com/kaleido-io/firefly/mocks/dataexchangemocks"
	"github.com/kaleido-io/firefly/pkg/fftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestDataManager(t *testing.T) (*dataManager, context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	mdi := &databasemocks.Plugin{}
	mdx := &dataexchangemocks.Plugin{}
	dm, err := NewDataManager(ctx, mdi, mdx)
	assert.NoError(t, err)
	return dm.(*dataManager), ctx, cancel
}

func TestValidateE2E(t *testing.T) {

	config.Reset()
	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)
	data := &fftypes.Data{
		Namespace: "ns1",
		Validator: fftypes.ValidatorTypeJSON,
		Datatype: &fftypes.DatatypeRef{
			Name:    "customer",
			Version: "0.0.1",
		},
		Value: fftypes.Byteable(`{"some":"json"}`),
	}
	data.Seal(ctx)
	dt := &fftypes.Datatype{
		ID:        fftypes.NewUUID(),
		Validator: fftypes.ValidatorTypeJSON,
		Value: fftypes.Byteable(`{
			"properties": {
				"field1": {
					"type": "string"
				}
			},
			"additionalProperties": false
		}`),
		Namespace: "ns1",
		Name:      "customer",
		Version:   "0.0.1",
	}
	mdi.On("GetDatatypeByName", mock.Anything, "ns1", "customer", "0.0.1").Return(dt, nil)
	isValid, err := dm.ValidateAll(ctx, []*fftypes.Data{data})
	assert.Regexp(t, "FF10198", err)
	assert.False(t, isValid)

	v, err := dm.getValidatorForDatatype(ctx, data.Namespace, data.Validator, data.Datatype)
	err = v.Validate(ctx, data)
	assert.Regexp(t, "FF10198", err)

	data.Value = fftypes.Byteable(`{"field1":"value1"}`)
	data.Seal(context.Background())
	err = v.Validate(ctx, data)
	assert.NoError(t, err)

	isValid, err = dm.ValidateAll(ctx, []*fftypes.Data{data})
	assert.NoError(t, err)
	assert.True(t, isValid)

}

func TestInitBadDeps(t *testing.T) {
	_, err := NewDataManager(context.Background(), nil, nil)
	assert.Regexp(t, "FF10128", err)
}

func TestValidatorLookupCached(t *testing.T) {

	config.Reset()
	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)
	ref := &fftypes.DatatypeRef{
		Name:    "customer",
		Version: "0.0.1",
	}
	dt := &fftypes.Datatype{
		ID:        fftypes.NewUUID(),
		Validator: fftypes.ValidatorTypeJSON,
		Value:     fftypes.Byteable(`{}`),
		Name:      "customer",
		Namespace: "0.0.1",
	}
	mdi.On("GetDatatypeByName", mock.Anything, "ns1", "customer", "0.0.1").Return(dt, nil).Once()
	lookup1, err := dm.getValidatorForDatatype(ctx, "ns1", fftypes.ValidatorTypeJSON, ref)
	assert.NoError(t, err)
	assert.Equal(t, "customer", lookup1.(*jsonValidator).datatype.Name)

	lookup2, err := dm.getValidatorForDatatype(ctx, "ns1", fftypes.ValidatorTypeJSON, ref)
	assert.NoError(t, err)
	assert.Equal(t, lookup1, lookup2)

}

func TestValidateBadHash(t *testing.T) {

	config.Reset()
	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)
	data := &fftypes.Data{
		Namespace: "ns1",
		Validator: fftypes.ValidatorTypeJSON,
		Datatype: &fftypes.DatatypeRef{
			Name:    "customer",
			Version: "0.0.1",
		},
		Value: fftypes.Byteable(`{}`),
		Hash:  fftypes.NewRandB32(),
	}
	dt := &fftypes.Datatype{
		ID:        fftypes.NewUUID(),
		Validator: fftypes.ValidatorTypeJSON,
		Value:     fftypes.Byteable(`{}`),
		Name:      "customer",
		Namespace: "0.0.1",
	}
	mdi.On("GetDatatypeByName", mock.Anything, "ns1", "customer", "0.0.1").Return(dt, nil).Once()
	_, err := dm.ValidateAll(ctx, []*fftypes.Data{data})
	assert.Regexp(t, "FF10201", err)

}

func TestGetMessageDataDBError(t *testing.T) {

	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)
	mdi.On("GetDataByID", mock.Anything, mock.Anything, true).Return(nil, fmt.Errorf("pop"))
	data, foundAll, err := dm.GetMessageData(ctx, &fftypes.Message{
		Header: fftypes.MessageHeader{ID: fftypes.NewUUID()},
		Data:   fftypes.DataRefs{{ID: fftypes.NewUUID(), Hash: fftypes.NewRandB32()}},
	}, true)
	assert.Nil(t, data)
	assert.False(t, foundAll)
	assert.EqualError(t, err, "pop")

}

func TestGetMessageDataNilEntry(t *testing.T) {

	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)
	mdi.On("GetDataByID", mock.Anything, mock.Anything, true).Return(nil, nil)
	data, foundAll, err := dm.GetMessageData(ctx, &fftypes.Message{
		Header: fftypes.MessageHeader{ID: fftypes.NewUUID()},
		Data:   fftypes.DataRefs{nil},
	}, true)
	assert.Empty(t, data)
	assert.False(t, foundAll)
	assert.NoError(t, err)

}

func TestGetMessageDataNotFound(t *testing.T) {

	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)
	mdi.On("GetDataByID", mock.Anything, mock.Anything, true).Return(nil, nil)
	data, foundAll, err := dm.GetMessageData(ctx, &fftypes.Message{
		Header: fftypes.MessageHeader{ID: fftypes.NewUUID()},
		Data:   fftypes.DataRefs{{ID: fftypes.NewUUID(), Hash: fftypes.NewRandB32()}},
	}, true)
	assert.Empty(t, data)
	assert.False(t, foundAll)
	assert.NoError(t, err)

}

func TestGetMessageDataHashMismatch(t *testing.T) {

	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)
	dataID := fftypes.NewUUID()
	mdi.On("GetDataByID", mock.Anything, mock.Anything, true).Return(&fftypes.Data{
		ID:   dataID,
		Hash: fftypes.NewRandB32(),
	}, nil)
	data, foundAll, err := dm.GetMessageData(ctx, &fftypes.Message{
		Header: fftypes.MessageHeader{ID: fftypes.NewUUID()},
		Data:   fftypes.DataRefs{{ID: dataID, Hash: fftypes.NewRandB32()}},
	}, true)
	assert.Empty(t, data)
	assert.False(t, foundAll)
	assert.NoError(t, err)

}

func TestGetMessageDataOk(t *testing.T) {

	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)
	dataID := fftypes.NewUUID()
	hash := fftypes.NewRandB32()
	mdi.On("GetDataByID", mock.Anything, mock.Anything, true).Return(&fftypes.Data{
		ID:   dataID,
		Hash: hash,
	}, nil)
	data, foundAll, err := dm.GetMessageData(ctx, &fftypes.Message{
		Header: fftypes.MessageHeader{ID: fftypes.NewUUID()},
		Data:   fftypes.DataRefs{{ID: dataID, Hash: hash}},
	}, true)
	assert.NotEmpty(t, data)
	assert.Equal(t, *dataID, *data[0].ID)
	assert.True(t, foundAll)
	assert.NoError(t, err)

}

func TestCheckDatatypeVerifiesTheSchema(t *testing.T) {

	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	err := dm.CheckDatatype(ctx, "ns1", &fftypes.Datatype{})
	assert.Regexp(t, "FF10196", err)
}

func TestResolveInputDataEmpty(t *testing.T) {

	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	refs, err := dm.ResolveInputData(ctx, "ns1", fftypes.InputData{})
	assert.NoError(t, err)
	assert.Empty(t, refs)

}

func TestResolveInputDataRefIDOnlyOK(t *testing.T) {
	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)

	dataID := fftypes.NewUUID()
	dataHash := fftypes.NewRandB32()

	mdi.On("GetDataByID", ctx, dataID, false).Return(&fftypes.Data{
		ID:        dataID,
		Namespace: "ns1",
		Hash:      dataHash,
	}, nil)

	refs, err := dm.ResolveInputData(ctx, "ns1", fftypes.InputData{
		{DataRef: fftypes.DataRef{ID: dataID}},
	})
	assert.NoError(t, err)
	assert.Len(t, refs, 1)
	assert.Equal(t, dataID, refs[0].ID)
	assert.Equal(t, dataHash, refs[0].Hash)
}

func TestResolveInputDataRefBadNamespace(t *testing.T) {
	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)

	dataID := fftypes.NewUUID()
	dataHash := fftypes.NewRandB32()

	mdi.On("GetDataByID", ctx, dataID, false).Return(&fftypes.Data{
		ID:        dataID,
		Namespace: "ns2",
		Hash:      dataHash,
	}, nil)

	refs, err := dm.ResolveInputData(ctx, "ns1", fftypes.InputData{
		{DataRef: fftypes.DataRef{ID: dataID, Hash: dataHash}},
	})
	assert.Regexp(t, "FF10204", err)
	assert.Empty(t, refs)
}

func TestResolveInputDataRefBadHash(t *testing.T) {
	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)

	dataID := fftypes.NewUUID()
	dataHash := fftypes.NewRandB32()

	mdi.On("GetDataByID", ctx, dataID, false).Return(&fftypes.Data{
		ID:        dataID,
		Namespace: "ns2",
		Hash:      dataHash,
	}, nil)

	refs, err := dm.ResolveInputData(ctx, "ns1", fftypes.InputData{
		{DataRef: fftypes.DataRef{ID: dataID, Hash: fftypes.NewRandB32()}},
	})
	assert.Regexp(t, "FF10204", err)
	assert.Empty(t, refs)
}

func TestResolveInputDataRefLookkupFail(t *testing.T) {
	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)

	dataID := fftypes.NewUUID()

	mdi.On("GetDataByID", ctx, dataID, false).Return(nil, fmt.Errorf("pop"))

	_, err := dm.ResolveInputData(ctx, "ns1", fftypes.InputData{
		{DataRef: fftypes.DataRef{ID: dataID, Hash: fftypes.NewRandB32()}},
	})
	assert.EqualError(t, err, "pop")
}

func TestResolveInputDataValueNoValidatorOK(t *testing.T) {
	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)

	mdi.On("UpsertData", ctx, mock.Anything, false, false).Return(nil)

	refs, err := dm.ResolveInputData(ctx, "ns1", fftypes.InputData{
		{Value: fftypes.Byteable(`{"some":"json"}`)},
	})
	assert.NoError(t, err)
	assert.Len(t, refs, 1)
	assert.NotNil(t, refs[0].ID)
	assert.NotNil(t, refs[0].Hash)
}

func TestResolveInputDataValueNoValidatorStoreFail(t *testing.T) {
	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)

	mdi.On("UpsertData", ctx, mock.Anything, false, false).Return(fmt.Errorf("pop"))

	_, err := dm.ResolveInputData(ctx, "ns1", fftypes.InputData{
		{Value: fftypes.Byteable(`{"some":"json"}`)},
	})
	assert.EqualError(t, err, "pop")
}

func TestResolveInputDataValueWithValidation(t *testing.T) {
	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)

	mdi.On("UpsertData", ctx, mock.Anything, false, false).Return(nil)
	mdi.On("GetDatatypeByName", ctx, "ns1", "customer", "0.0.1").Return(&fftypes.Datatype{
		ID:        fftypes.NewUUID(),
		Validator: fftypes.ValidatorTypeJSON,
		Namespace: "ns1",
		Name:      "customer",
		Version:   "0.0.1",
		Value: fftypes.Byteable(`{
			"properties": {
				"field1": {
					"type": "string"
				}
			},
			"additionalProperties": false
		}`),
	}, nil)

	refs, err := dm.ResolveInputData(ctx, "ns1", fftypes.InputData{
		{
			Datatype: &fftypes.DatatypeRef{
				Name:    "customer",
				Version: "0.0.1",
			},
			Value: fftypes.Byteable(`{"field1":"value1"}`),
		},
	})
	assert.NoError(t, err)
	assert.Len(t, refs, 1)
	assert.NotNil(t, refs[0].ID)
	assert.NotNil(t, refs[0].Hash)

	_, err = dm.ResolveInputData(ctx, "ns1", fftypes.InputData{
		{
			Datatype: &fftypes.DatatypeRef{
				Name:    "customer",
				Version: "0.0.1",
			},
			Value: fftypes.Byteable(`{"not_allowed":"value"}`),
		},
	})
	assert.Regexp(t, "FF10198", err)
}

func TestResolveInputDataNoRefOrValue(t *testing.T) {
	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()

	_, err := dm.ResolveInputData(ctx, "ns1", fftypes.InputData{
		{ /* missing */ },
	})
	assert.Regexp(t, "FF10205", err)
}

func TestUploadJSONLoadDatatypeFail(t *testing.T) {
	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)

	mdi.On("GetDatatypeByName", ctx, "ns1", "customer", "0.0.1").Return(nil, fmt.Errorf("pop"))
	_, err := dm.UploadJSON(ctx, "ns1", &fftypes.Data{
		Datatype: &fftypes.DatatypeRef{
			Name:    "customer",
			Version: "0.0.1",
		},
	})
	assert.EqualError(t, err, "pop")
}

func TestValidateAndStoreLoadNilRef(t *testing.T) {
	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()

	_, err := dm.validateAndStoreInlined(ctx, "ns1", &fftypes.DataRefOrValue{
		Validator: fftypes.ValidatorTypeJSON,
		Datatype:  nil,
	})
	assert.Regexp(t, "FF10199", err)
}

func TestValidateAndStoreLoadValidatorUnknown(t *testing.T) {

	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)
	mdi.On("GetDatatypeByName", mock.Anything, "ns1", "customer", "0.0.1").Return(nil, nil)
	_, err := dm.validateAndStoreInlined(ctx, "ns1", &fftypes.DataRefOrValue{
		Validator: "wrong!",
		Datatype: &fftypes.DatatypeRef{
			Name:    "customer",
			Version: "0.0.1",
		},
	})
	assert.Regexp(t, "FF10200.*wrong", err)

}

func TestValidateAndStoreLoadBadRef(t *testing.T) {

	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)
	mdi.On("GetDatatypeByName", mock.Anything, "ns1", "customer", "0.0.1").Return(nil, nil)
	_, err := dm.validateAndStoreInlined(ctx, "ns1", &fftypes.DataRefOrValue{
		Datatype: &fftypes.DatatypeRef{
			// Missing name
		},
	})
	assert.Regexp(t, "FF10195", err)
}

func TestValidateAndStoreNotFound(t *testing.T) {

	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)
	mdi.On("GetDatatypeByName", mock.Anything, "ns1", "customer", "0.0.1").Return(nil, nil)
	_, err := dm.validateAndStoreInlined(ctx, "ns1", &fftypes.DataRefOrValue{
		Datatype: &fftypes.DatatypeRef{
			Name:    "customer",
			Version: "0.0.1",
		},
	})
	assert.Regexp(t, "FF10195", err)
}

func TestValidateAllLookupError(t *testing.T) {

	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)
	mdi.On("GetDatatypeByName", mock.Anything, "ns1", "customer", "0.0.1").Return(nil, fmt.Errorf("pop"))
	data := &fftypes.Data{
		Namespace: "ns1",
		Validator: fftypes.ValidatorTypeJSON,
		Datatype: &fftypes.DatatypeRef{
			Name:    "customer",
			Version: "0.0.1",
		},
		Value: fftypes.Byteable(`anything`),
	}
	data.Seal(ctx)
	_, err := dm.ValidateAll(ctx, []*fftypes.Data{data})
	assert.Regexp(t, "pop", err)

}

func TestGetValidatorForDatatypeNilRef(t *testing.T) {

	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	v, err := dm.getValidatorForDatatype(ctx, "", "", nil)
	assert.Nil(t, v)
	assert.NoError(t, err)

}

func TestValidateAllStoredValidatorInvalid(t *testing.T) {

	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)
	mdi.On("GetDatatypeByName", mock.Anything, "ns1", "customer", "0.0.1").Return(&fftypes.Datatype{
		Value: fftypes.Byteable(`{"not": "a", "schema": true}`),
	}, nil)
	data := &fftypes.Data{
		Namespace: "ns1",
		Datatype: &fftypes.DatatypeRef{
			Name:    "customer",
			Version: "0.0.1",
		},
	}
	isValid, err := dm.ValidateAll(ctx, []*fftypes.Data{data})
	assert.False(t, isValid)
	assert.NoError(t, err)
	mdi.AssertExpectations(t)
}

func TestVerifyNamespaceExistsInvalidFFName(t *testing.T) {
	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	err := dm.VerifyNamespaceExists(ctx, "!wrong")
	assert.Regexp(t, "FF10131", err)
}

func TestVerifyNamespaceExistsLookupErr(t *testing.T) {
	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)
	mdi.On("GetNamespace", mock.Anything, "ns1").Return(nil, fmt.Errorf("pop"))
	err := dm.VerifyNamespaceExists(ctx, "ns1")
	assert.Regexp(t, "pop", err)
}

func TestVerifyNamespaceExistsNotFound(t *testing.T) {
	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)
	mdi.On("GetNamespace", mock.Anything, "ns1").Return(nil, nil)
	err := dm.VerifyNamespaceExists(ctx, "ns1")
	assert.Regexp(t, "FF10187", err)
}

func TestVerifyNamespaceExistsOk(t *testing.T) {
	dm, ctx, cancel := newTestDataManager(t)
	defer cancel()
	mdi := dm.database.(*databasemocks.Plugin)
	mdi.On("GetNamespace", mock.Anything, "ns1").Return(&fftypes.Namespace{}, nil)
	err := dm.VerifyNamespaceExists(ctx, "ns1")
	assert.NoError(t, err)
}
