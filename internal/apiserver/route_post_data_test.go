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

package apiserver

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/hyperledger-labs/firefly/internal/log"
	"github.com/hyperledger-labs/firefly/mocks/datamocks"
	"github.com/hyperledger-labs/firefly/mocks/orchestratormocks"
	"github.com/hyperledger-labs/firefly/pkg/fftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostDataJSON(t *testing.T) {
	o := &orchestratormocks.Orchestrator{}
	mdm := &datamocks.Manager{}
	o.On("Data").Return(mdm)
	as := &apiServer{}
	r := as.createMuxRouter(o)
	input := fftypes.Data{}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(&input)
	req := httptest.NewRequest("POST", "/api/v1/namespaces/ns1/data", &buf)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	res := httptest.NewRecorder()

	mdm.On("UploadJSON", mock.Anything, "ns1", mock.AnythingOfType("*fftypes.Data")).
		Return(&fftypes.Data{}, nil)
	r.ServeHTTP(res, req)

	assert.Equal(t, 201, res.Result().StatusCode)
}

func TestPostDataBinary(t *testing.T) {
	log.SetLevel("debug")

	o := &orchestratormocks.Orchestrator{}
	mdm := &datamocks.Manager{}
	o.On("Data").Return(mdm)
	as := &apiServer{}
	r := as.createMuxRouter(o)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	writer, err := w.CreateFormFile("file", "filename.ext")
	assert.NoError(t, err)
	writer.Write([]byte(`some data`))
	w.Close()
	req := httptest.NewRequest("POST", "/api/v1/namespaces/ns1/data", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	res := httptest.NewRecorder()

	mdm.On("UploadBLOB", mock.Anything, "ns1", mock.AnythingOfType("*fftypes.Data"), mock.AnythingOfType("*fftypes.Multipart"), false).
		Return(&fftypes.Data{}, nil)
	r.ServeHTTP(res, req)

	assert.Equal(t, 201, res.Result().StatusCode)
}

func TestPostDataBinaryObjAutoMeta(t *testing.T) {
	log.SetLevel("debug")

	o := &orchestratormocks.Orchestrator{}
	mdm := &datamocks.Manager{}
	o.On("Data").Return(mdm)
	as := &apiServer{}
	r := as.createMuxRouter(o)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	writer, err := w.CreateFormField("metadata")
	assert.NoError(t, err)
	writer.Write([]byte(`{"filename":"anything"}`))
	writer, err = w.CreateFormField("validator")
	assert.NoError(t, err)
	writer.Write([]byte(fftypes.ValidatorTypeJSON))
	writer, err = w.CreateFormField("datatype.name")
	assert.NoError(t, err)
	writer.Write([]byte("fileinfo"))
	writer, err = w.CreateFormField("datatype.version")
	assert.NoError(t, err)
	writer.Write([]byte("0.0.1"))
	writer, err = w.CreateFormField("autometa")
	assert.NoError(t, err)
	writer.Write([]byte("true"))
	writer, err = w.CreateFormFile("file", "filename.ext")
	assert.NoError(t, err)
	writer.Write([]byte(`some data`))
	w.Close()
	req := httptest.NewRequest("POST", "/api/v1/namespaces/ns1/data", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	res := httptest.NewRecorder()

	mdm.On("UploadBLOB", mock.Anything, "ns1", mock.MatchedBy(func(d *fftypes.Data) bool {
		assert.Equal(t, `{"filename":"anything"}`, string(d.Value))
		assert.Equal(t, fftypes.ValidatorTypeJSON, d.Validator)
		assert.Equal(t, "fileinfo", d.Datatype.Name)
		assert.Equal(t, "0.0.1", d.Datatype.Version)
		return true
	}), mock.AnythingOfType("*fftypes.Multipart"), true).
		Return(&fftypes.Data{}, nil)
	r.ServeHTTP(res, req)

	assert.Equal(t, 201, res.Result().StatusCode)
}

func TestPostDataBinaryStringMetadata(t *testing.T) {
	log.SetLevel("debug")

	o := &orchestratormocks.Orchestrator{}
	mdm := &datamocks.Manager{}
	o.On("Data").Return(mdm)
	as := &apiServer{}
	r := as.createMuxRouter(o)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	writer, err := w.CreateFormField("metadata")
	assert.NoError(t, err)
	writer.Write([]byte(`string metadata`))
	writer, err = w.CreateFormFile("file", "filename.ext")
	assert.NoError(t, err)
	writer.Write([]byte(`some data`))
	w.Close()
	req := httptest.NewRequest("POST", "/api/v1/namespaces/ns1/data", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	res := httptest.NewRecorder()

	mdm.On("UploadBLOB", mock.Anything, "ns1", mock.MatchedBy(func(d *fftypes.Data) bool {
		assert.Equal(t, `"string metadata"`, string(d.Value))
		assert.Equal(t, "", string(d.Validator))
		assert.Nil(t, d.Datatype)
		return true
	}), mock.AnythingOfType("*fftypes.Multipart"), false).
		Return(&fftypes.Data{}, nil)
	r.ServeHTTP(res, req)

	assert.Equal(t, 201, res.Result().StatusCode)
}

func TestPostDataTrailingMetadata(t *testing.T) {
	log.SetLevel("debug")

	o := &orchestratormocks.Orchestrator{}
	mdm := &datamocks.Manager{}
	o.On("Data").Return(mdm)
	as := &apiServer{}
	r := as.createMuxRouter(o)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	writer, err := w.CreateFormFile("file", "filename.ext")
	assert.NoError(t, err)
	writer.Write([]byte(`some data`))
	writer, err = w.CreateFormField("metadata")
	assert.NoError(t, err)
	writer.Write([]byte(`string metadata`))
	w.Close()
	req := httptest.NewRequest("POST", "/api/v1/namespaces/ns1/data", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	res := httptest.NewRecorder()

	mdm.On("UploadBLOB", mock.Anything, "ns1", mock.Anything, mock.AnythingOfType("*fftypes.Multipart"), false).
		Return(&fftypes.Data{}, nil)
	r.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Result().StatusCode)
	d, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Regexp(t, "FF10236.*metadata", string(d))
}

func TestPostDataBinaryMissing(t *testing.T) {
	log.SetLevel("debug")

	o := &orchestratormocks.Orchestrator{}
	mdm := &datamocks.Manager{}
	o.On("Data").Return(mdm)
	as := &apiServer{}
	r := as.createMuxRouter(o)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	formField, _ := w.CreateFormField("Not a file")
	formField.Write([]byte(`some value`))
	w.Close()
	req := httptest.NewRequest("POST", "/api/v1/namespaces/ns1/data", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	res := httptest.NewRecorder()

	mdm.On("UploadBLOB", mock.Anything, "ns1", mock.AnythingOfType("*fftypes.Data"), mock.AnythingOfType("*fftypes.Multipart"), false).
		Return(&fftypes.Data{}, nil)
	r.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Result().StatusCode)
}

func TestPostDataBadForm(t *testing.T) {
	log.SetLevel("debug")

	o := &orchestratormocks.Orchestrator{}
	mdm := &datamocks.Manager{}
	o.On("Data").Return(mdm)
	as := &apiServer{}
	r := as.createMuxRouter(o)

	var b bytes.Buffer
	req := httptest.NewRequest("POST", "/api/v1/namespaces/ns1/data", &b)
	req.Header.Set("Content-Type", "multipart/form-data")

	res := httptest.NewRecorder()

	mdm.On("UploadBLOB", mock.Anything, "ns1", mock.AnythingOfType("*fftypes.Data"), mock.AnythingOfType("*fftypes.Multipart"), false).
		Return(&fftypes.Data{}, nil)
	r.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Result().StatusCode)
}
