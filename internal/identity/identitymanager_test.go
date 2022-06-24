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

package identity

import (
	"context"
	"fmt"
	"testing"

	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly/internal/coreconfig"
	"github.com/hyperledger/firefly/mocks/blockchainmocks"
	"github.com/hyperledger/firefly/mocks/databasemocks"
	"github.com/hyperledger/firefly/mocks/multipartymocks"
	"github.com/hyperledger/firefly/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestIdentityManager(t *testing.T) (context.Context, *identityManager) {
	coreconfig.Reset()

	mdi := &databasemocks.Plugin{}
	mbi := &blockchainmocks.Plugin{}
	mmp := &multipartymocks.Manager{}

	mbi.On("VerifierType").Return(core.VerifierTypeEthAddress).Maybe()

	ctx := context.Background()
	im, err := NewIdentityManager(ctx, "ns1", "", "", "", mdi, mbi, mmp)
	assert.NoError(t, err)
	return ctx, im.(*identityManager)
}

func TestNewIdentityManagerMissingDeps(t *testing.T) {
	_, err := NewIdentityManager(context.Background(), "", "", "", "", nil, nil, nil)
	assert.Regexp(t, "FF10128", err)
}

func TestResolveInputSigningIdentityNoKey(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	msgIdentity := &core.SignerRef{}
	err := im.ResolveInputSigningIdentity(ctx, msgIdentity)
	assert.Regexp(t, "FF10354", err)

}

func TestResolveInputSigningIdentityOrgFallbackOk(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	im.orgName = "org1"
	im.orgKey = "key123"

	mbi := im.blockchain.(*blockchainmocks.Plugin)
	mbi.On("NormalizeSigningKey", ctx, "key123").Return("fullkey123", nil)

	orgID := fftypes.NewUUID()

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, "ns1", "fullkey123").
		Return((&core.Verifier{
			Identity:  orgID,
			Namespace: "ns1",
			VerifierRef: core.VerifierRef{
				Type:  core.VerifierTypeEthAddress,
				Value: "fullkey123",
			},
		}).Seal(), nil)
	mdi.On("GetIdentityByID", ctx, "ns1", orgID).
		Return(&core.Identity{
			IdentityBase: core.IdentityBase{
				ID:        orgID,
				DID:       "did:firefly:org/org1",
				Namespace: "ns1",
				Name:      "org1",
				Type:      core.IdentityTypeOrg,
			},
		}, nil)

	msgIdentity := &core.SignerRef{}
	err := im.ResolveInputSigningIdentity(ctx, msgIdentity)
	assert.NoError(t, err)
	assert.Equal(t, "did:firefly:org/org1", msgIdentity.Author)
	assert.Equal(t, "fullkey123", msgIdentity.Key)

	mbi.AssertExpectations(t)
	mdi.AssertExpectations(t)

}

func TestResolveInputSigningIdentityByKeyOk(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	mbi := im.blockchain.(*blockchainmocks.Plugin)
	mbi.On("NormalizeSigningKey", ctx, "mykey123").Return("fullkey123", nil)

	idID := fftypes.NewUUID()

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, "ns1", "fullkey123").
		Return((&core.Verifier{
			Identity:  idID,
			Namespace: "ns1",
			VerifierRef: core.VerifierRef{
				Type:  core.VerifierTypeEthAddress,
				Value: "fullkey123",
			},
		}).Seal(), nil)
	mdi.On("GetIdentityByID", ctx, "ns1", idID).
		Return(&core.Identity{
			IdentityBase: core.IdentityBase{
				ID:        idID,
				DID:       "did:firefly:ns/ns1/myid",
				Namespace: "ns1",
				Name:      "myid",
				Type:      core.IdentityTypeCustom,
			},
		}, nil)

	msgIdentity := &core.SignerRef{
		Key: "mykey123",
	}
	err := im.ResolveInputSigningIdentity(ctx, msgIdentity)
	assert.NoError(t, err)
	assert.Equal(t, "did:firefly:ns/ns1/myid", msgIdentity.Author)
	assert.Equal(t, "fullkey123", msgIdentity.Key)

	mbi.AssertExpectations(t)
	mdi.AssertExpectations(t)

}

func TestResolveInputSigningIdentityAnonymousKeyWithAuthorOk(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	mbi := im.blockchain.(*blockchainmocks.Plugin)
	mmp := im.multiparty.(*multipartymocks.Manager)
	mbi.On("NormalizeSigningKey", ctx, "mykey123").Return("fullkey123", nil)
	mmp.On("GetNetworkVersion").Return(1)

	idID := fftypes.NewUUID()

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, "ns1", "fullkey123").Return(nil, nil)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, "ns1", "fullkey123").Return(nil, nil)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, core.LegacySystemNamespace, "fullkey123").Return(nil, nil)
	mdi.On("GetIdentityByDID", ctx, "ns1", "did:firefly:ns/ns1/myid").
		Return(&core.Identity{
			IdentityBase: core.IdentityBase{
				ID:        idID,
				DID:       "did:firefly:ns/ns1/myid",
				Namespace: "ns1",
				Name:      "myid",
				Type:      core.IdentityTypeCustom,
			},
		}, nil)

	msgIdentity := &core.SignerRef{
		Key:    "mykey123",
		Author: "did:firefly:ns/ns1/myid",
	}
	err := im.ResolveInputSigningIdentity(ctx, msgIdentity)
	assert.NoError(t, err)
	assert.Equal(t, "did:firefly:ns/ns1/myid", msgIdentity.Author)
	assert.Equal(t, "fullkey123", msgIdentity.Key)

	mbi.AssertExpectations(t)
	mdi.AssertExpectations(t)

}

func TestResolveInputSigningIdentityKeyWithNoAuthorFail(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	mbi := im.blockchain.(*blockchainmocks.Plugin)
	mmp := im.multiparty.(*multipartymocks.Manager)
	mbi.On("NormalizeSigningKey", ctx, "mykey123").Return("fullkey123", nil)
	mmp.On("GetNetworkVersion").Return(1)

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, "ns1", "fullkey123").Return(nil, nil)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, "ns1", "fullkey123").Return(nil, nil)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, core.LegacySystemNamespace, "fullkey123").Return(nil, nil)

	msgIdentity := &core.SignerRef{
		Key: "mykey123",
	}
	err := im.ResolveInputSigningIdentity(ctx, msgIdentity)
	assert.Regexp(t, "FF10356", err)

	mbi.AssertExpectations(t)
	mdi.AssertExpectations(t)
	mmp.AssertExpectations(t)
}

func TestResolveInputSigningIdentityByKeyDIDMismatch(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	mbi := im.blockchain.(*blockchainmocks.Plugin)
	mbi.On("NormalizeSigningKey", ctx, "mykey123").Return("fullkey123", nil)

	idID := fftypes.NewUUID()

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, "ns1", "fullkey123").
		Return((&core.Verifier{
			Identity:  idID,
			Namespace: "ns1",
			VerifierRef: core.VerifierRef{
				Type:  core.VerifierTypeEthAddress,
				Value: "fullkey123",
			},
		}).Seal(), nil)
	mdi.On("GetIdentityByID", ctx, "ns1", idID).
		Return(&core.Identity{
			IdentityBase: core.IdentityBase{
				ID:        idID,
				DID:       "did:firefly:ns/ns1/myid",
				Namespace: "ns1",
				Name:      "myid",
				Type:      core.IdentityTypeCustom,
			},
		}, nil)

	msgIdentity := &core.SignerRef{
		Key:    "mykey123",
		Author: "did:firefly:ns/ns1/notmyid",
	}
	err := im.ResolveInputSigningIdentity(ctx, msgIdentity)
	assert.Regexp(t, "FF10355", err)

	mbi.AssertExpectations(t)
	mdi.AssertExpectations(t)

}

func TestResolveInputSigningIdentityByKeyNotFound(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	mbi := im.blockchain.(*blockchainmocks.Plugin)
	mmp := im.multiparty.(*multipartymocks.Manager)
	mbi.On("NormalizeSigningKey", ctx, "mykey123").Return("fullkey123", nil)
	mmp.On("GetNetworkVersion").Return(1)

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, "ns1", "fullkey123").
		Return(nil, nil)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, "ns1", "fullkey123").
		Return(nil, nil)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, core.LegacySystemNamespace, "fullkey123").
		Return(nil, nil)
	mdi.On("GetIdentityByDID", ctx, "ns1", "did:firefly:ns/ns1/unknown").
		Return(nil, nil)

	msgIdentity := &core.SignerRef{
		Key:    "mykey123",
		Author: "did:firefly:ns/ns1/unknown",
	}
	err := im.ResolveInputSigningIdentity(ctx, msgIdentity)
	assert.Regexp(t, "FF10277", err)

	mbi.AssertExpectations(t)
	mdi.AssertExpectations(t)

}

func TestResolveInputSigningIdentityByKeyFail(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	mbi := im.blockchain.(*blockchainmocks.Plugin)
	mbi.On("NormalizeSigningKey", ctx, "mykey123").Return("fullkey123", nil)

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, "ns1", "fullkey123").
		Return(nil, fmt.Errorf("pop"))

	msgIdentity := &core.SignerRef{
		Key: "mykey123",
	}
	err := im.ResolveInputSigningIdentity(ctx, msgIdentity)
	assert.Regexp(t, "pop", err)

	mbi.AssertExpectations(t)
	mdi.AssertExpectations(t)

}

func TestResolveInputSigningIdentityByKeyResolveFail(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	mbi := im.blockchain.(*blockchainmocks.Plugin)
	mbi.On("NormalizeSigningKey", ctx, "mykey123").Return("", fmt.Errorf("pop"))

	msgIdentity := &core.SignerRef{
		Key: "mykey123",
	}
	err := im.ResolveInputSigningIdentity(ctx, msgIdentity)
	assert.Regexp(t, "pop", err)

	mbi.AssertExpectations(t)
}

func TestResolveInputSigningIdentityByOrgNameOk(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	idID := fftypes.NewUUID()

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetIdentityByName", ctx, core.IdentityTypeOrg, "ns1", "org1").
		Return(&core.Identity{
			IdentityBase: core.IdentityBase{
				ID:        idID,
				DID:       "did:firefly:org/org1",
				Namespace: "ns1",
				Name:      "myid",
				Type:      core.IdentityTypeOrg,
			},
		}, nil)
	mdi.On("GetVerifiers", ctx, "ns1", mock.Anything).
		Return([]*core.Verifier{
			(&core.Verifier{
				Identity:  idID,
				Namespace: "ns1",
				VerifierRef: core.VerifierRef{
					Type:  core.VerifierTypeEthAddress,
					Value: "fullkey123",
				},
			}).Seal(),
		}, nil, nil)

	msgIdentity := &core.SignerRef{
		Author: "org1",
	}
	err := im.ResolveInputSigningIdentity(ctx, msgIdentity)
	assert.NoError(t, err)
	assert.Equal(t, "did:firefly:org/org1", msgIdentity.Author)
	assert.Equal(t, "fullkey123", msgIdentity.Key)

	mdi.AssertExpectations(t)

}

func TestResolveInputSigningIdentityByOrgLookkupNotFound(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetIdentityByName", ctx, core.IdentityTypeOrg, "ns1", "org1").
		Return(nil, nil)

	msgIdentity := &core.SignerRef{
		Author: "org1",
	}
	err := im.ResolveInputSigningIdentity(ctx, msgIdentity)
	assert.Regexp(t, "FF10277", err)

	mdi.AssertExpectations(t)

}

func TestResolveInputSigningIdentityByOrgLookkupFail(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetIdentityByName", ctx, core.IdentityTypeOrg, "ns1", "org1").
		Return(nil, fmt.Errorf("pop"))

	msgIdentity := &core.SignerRef{
		Author: "org1",
	}
	err := im.ResolveInputSigningIdentity(ctx, msgIdentity)
	assert.Regexp(t, "pop", err)

	mdi.AssertExpectations(t)

}

func TestResolveInputSigningIdentityByOrgVerifierFail(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	idID := fftypes.NewUUID()

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetIdentityByName", ctx, core.IdentityTypeOrg, "ns1", "org1").
		Return(&core.Identity{
			IdentityBase: core.IdentityBase{
				ID:        idID,
				DID:       "did:firefly:org/org1",
				Namespace: "ns1",
				Name:      "myid",
				Type:      core.IdentityTypeOrg,
			},
		}, nil)
	mdi.On("GetVerifiers", ctx, "ns1", mock.Anything).Return(nil, nil, fmt.Errorf("pop"))

	msgIdentity := &core.SignerRef{
		Author: "org1",
	}
	err := im.ResolveInputSigningIdentity(ctx, msgIdentity)
	assert.Regexp(t, "pop", err)

	mdi.AssertExpectations(t)

}

func TestNormalizeSigningKeyDefault(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	im.defaultKey = "key123"

	mbi := im.blockchain.(*blockchainmocks.Plugin)
	mbi.On("NormalizeSigningKey", ctx, "key123").Return("fullkey123", nil)

	resolvedKey, err := im.NormalizeSigningKey(ctx, "", KeyNormalizationBlockchainPlugin)
	assert.NoError(t, err)
	assert.Equal(t, "fullkey123", resolvedKey)

	mbi.AssertExpectations(t)

}

func TestNormalizeSigningKeyOrgFallbackOk(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	im.orgKey = "key123"

	mbi := im.blockchain.(*blockchainmocks.Plugin)
	mbi.On("NormalizeSigningKey", ctx, "key123").Return("fullkey123", nil)

	resolvedKey, err := im.NormalizeSigningKey(ctx, "", KeyNormalizationBlockchainPlugin)
	assert.NoError(t, err)
	assert.Equal(t, "fullkey123", resolvedKey)

	mbi.AssertExpectations(t)

}

func TestNormalizeSigningKeyOrgFallbackErr(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	im.orgKey = "key123"

	mbi := im.blockchain.(*blockchainmocks.Plugin)
	mbi.On("NormalizeSigningKey", ctx, "key123").Return("fullkey123", fmt.Errorf("pop"))

	_, err := im.NormalizeSigningKey(ctx, "", KeyNormalizationBlockchainPlugin)
	assert.Regexp(t, "pop", err)

	mbi.AssertExpectations(t)

}

func TestResolveInputSigningKeyOk(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	mbi := im.blockchain.(*blockchainmocks.Plugin)
	mbi.On("NormalizeSigningKey", ctx, "key123").Return("fullkey123", nil)

	resolvedKey, err := im.NormalizeSigningKey(ctx, "key123", KeyNormalizationBlockchainPlugin)
	assert.NoError(t, err)
	assert.Equal(t, "fullkey123", resolvedKey)

	mbi.AssertExpectations(t)
}

func TestResolveInputSigningKeyFail(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	mbi := im.blockchain.(*blockchainmocks.Plugin)
	mbi.On("NormalizeSigningKey", ctx, "key123").Return("", fmt.Errorf("pop"))

	_, err := im.NormalizeSigningKey(ctx, "key123", KeyNormalizationBlockchainPlugin)
	assert.Regexp(t, "pop", err)

	mbi.AssertExpectations(t)
}

func TestResolveInputSigningKeyBypass(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	key, err := im.NormalizeSigningKey(ctx, "different-type-of-key", KeyNormalizationNone)
	assert.NoError(t, err)
	assert.Equal(t, "different-type-of-key", key)
}

func TestFirstVerifierForIdentityNotFound(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	id := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			DID:       "did:firefly:org/org1",
			Namespace: "ns1",
			Name:      "myid",
			Type:      core.IdentityTypeOrg,
		},
	}

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetVerifiers", ctx, "ns1", mock.Anything).Return([]*core.Verifier{}, nil, nil)

	_, retryable, err := im.firstVerifierForIdentity(ctx, core.VerifierTypeEthAddress, id)
	assert.Regexp(t, "FF10353", err)
	assert.False(t, retryable)

	mdi.AssertExpectations(t)

}

func TestResolveDefaultSigningIdentityNotFound(t *testing.T) {

	ctx, im := newTestIdentityManager(t)
	im.multipartyRootVerifier = &core.VerifierRef{
		Type:  core.VerifierTypeEthAddress,
		Value: "key12345",
	}

	mbi := im.blockchain.(*blockchainmocks.Plugin)
	mmp := im.multiparty.(*multipartymocks.Manager)
	mmp.On("GetNetworkVersion").Return(1)

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, "ns1", "key12345").Return(nil, nil)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, core.LegacySystemNamespace, "key12345").Return(nil, nil)

	err := im.resolveDefaultSigningIdentity(ctx, &core.SignerRef{})
	assert.Regexp(t, "FF10281", err)

	mbi.AssertExpectations(t)
	mdi.AssertExpectations(t)

}

func TestResolveDefaultSigningIdentitySystemFallback(t *testing.T) {

	ctx, im := newTestIdentityManager(t)
	im.multipartyRootVerifier = &core.VerifierRef{
		Type:  core.VerifierTypeEthAddress,
		Value: "key12345",
	}

	id := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			DID:       "did:firefly:org/org1",
			Namespace: "ns1",
			Name:      "org1",
			Type:      core.IdentityTypeOrg,
		},
	}
	verifier := &core.Verifier{
		Identity: id.ID,
		VerifierRef: core.VerifierRef{
			Value: "key12345",
		},
	}

	mbi := im.blockchain.(*blockchainmocks.Plugin)
	mmp := im.multiparty.(*multipartymocks.Manager)
	mmp.On("GetNetworkVersion").Return(1)

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, "ns1", "key12345").Return(nil, nil)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, core.LegacySystemNamespace, "key12345").Return(verifier, nil)
	mdi.On("GetIdentityByID", ctx, core.LegacySystemNamespace, id.ID).Return(id, nil)

	im.orgName = "org1"

	ref := &core.SignerRef{}
	err := im.resolveDefaultSigningIdentity(ctx, ref)
	assert.NoError(t, err)
	assert.Equal(t, "did:firefly:org/org1", ref.Author)
	assert.Equal(t, "key12345", ref.Key)

	mbi.AssertExpectations(t)
	mdi.AssertExpectations(t)

}

func TestGetMultipartyRootVerifierResolveFailed(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	mbi := im.blockchain.(*blockchainmocks.Plugin)
	mbi.On("NormalizeSigningKey", ctx, "0x12345").Return("", fmt.Errorf("pop"))

	im.orgKey = "0x12345"

	_, err := im.GetMultipartyRootVerifier(ctx)
	assert.Regexp(t, "pop", err)

	mbi.AssertExpectations(t)

}

func TestNormalizeKeyViaBlockchainPluginEmptyRequest(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	_, err := im.normalizeKeyViaBlockchainPlugin(ctx, "")
	assert.Regexp(t, "FF10352", err)

}

func TestNormalizeKeyViaBlockchainPluginCached(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	mbi := im.blockchain.(*blockchainmocks.Plugin)
	mbi.On("NormalizeSigningKey", ctx, "0x12345").Return("resolved12345", nil).Once()

	v, err := im.normalizeKeyViaBlockchainPlugin(ctx, "0x12345")
	assert.NoError(t, err)
	assert.Equal(t, core.VerifierRef{
		Type:  core.VerifierTypeEthAddress,
		Value: "resolved12345",
	}, *v)

	v1, err := im.normalizeKeyViaBlockchainPlugin(ctx, "0x12345")
	assert.NoError(t, err)
	assert.Equal(t, v, v1)

}

func TestGetMultipartyRootOrgCached(t *testing.T) {

	ctx, im := newTestIdentityManager(t)
	im.multipartyRootOrg = &core.Identity{}

	id, err := im.GetMultipartyRootOrg(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, id)

}

func TestGetMultipartyRootVerifierNotSet(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	_, err := im.GetMultipartyRootOrg(ctx)
	assert.Regexp(t, "FF10354", err)

}

func TestGetMultipartyRootOrgMismatch(t *testing.T) {

	ctx, im := newTestIdentityManager(t)
	im.multipartyRootVerifier = &core.VerifierRef{
		Type:  core.VerifierTypeEthAddress,
		Value: "fullkey123",
	}

	orgID := fftypes.NewUUID()
	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, "ns1", "fullkey123").
		Return((&core.Verifier{
			Identity:  orgID,
			Namespace: "ns1",
			VerifierRef: core.VerifierRef{
				Type:  core.VerifierTypeEthAddress,
				Value: "fullkey123",
			},
		}).Seal(), nil)
	mdi.On("GetIdentityByID", ctx, "ns1", orgID).
		Return(&core.Identity{
			IdentityBase: core.IdentityBase{
				ID:        orgID,
				DID:       "did:firefly:org/org2",
				Namespace: "ns1",
				Name:      "org2",
				Type:      core.IdentityTypeOrg,
			},
		}, nil)

	_, err := im.GetMultipartyRootOrg(ctx)
	assert.Regexp(t, "FF10281", err)

}

func TestCachedIdentityLookupByVerifierRefCaching(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	id := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			DID:       "did:firefly:node/peer1",
			Namespace: "ns1",
			Name:      "peer1",
			Type:      core.IdentityTypeOrg,
		},
	}
	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeFFDXPeerID, "ns1", "peer1").
		Return((&core.Verifier{
			Identity:  id.ID,
			Namespace: "ns1",
			VerifierRef: core.VerifierRef{
				Type:  core.VerifierTypeFFDXPeerID,
				Value: "peer1",
			},
		}).Seal(), nil)
	mdi.On("GetIdentityByID", ctx, "ns1", id.ID).
		Return(id, nil)

	v1, err := im.cachedIdentityLookupByVerifierRef(ctx, "ns1", &core.VerifierRef{
		Type:  core.VerifierTypeFFDXPeerID,
		Value: "peer1",
	})
	assert.NoError(t, err)
	assert.Equal(t, id, v1)

	v2, err := im.cachedIdentityLookupByVerifierRef(ctx, "ns1", &core.VerifierRef{
		Type:  core.VerifierTypeFFDXPeerID,
		Value: "peer1",
	})
	assert.NoError(t, err)
	assert.Equal(t, id, v2)

}

func TestCachedIdentityLookupByVerifierRefError(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	id := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			DID:       "did:firefly:node/peer1",
			Namespace: "ns1",
			Name:      "peer1",
			Type:      core.IdentityTypeOrg,
		},
	}
	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, "ns1", "peer1").
		Return((&core.Verifier{
			Identity:  id.ID,
			Namespace: "ns1",
			VerifierRef: core.VerifierRef{
				Type:  core.VerifierTypeEthAddress,
				Value: "peer1",
			},
		}).Seal(), nil)
	mdi.On("GetIdentityByID", ctx, "ns1", id.ID).Return(nil, fmt.Errorf("pop"))

	_, err := im.cachedIdentityLookupByVerifierRef(ctx, "ns1", &core.VerifierRef{
		Type:  core.VerifierTypeEthAddress,
		Value: "peer1",
	})
	assert.Regexp(t, "pop", err)

}

func TestCachedIdentityLookupByVerifierRefNotFound(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	id := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			DID:       "did:firefly:node/peer1",
			Namespace: "ns1",
			Name:      "peer1",
			Type:      core.IdentityTypeOrg,
		},
	}
	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeEthAddress, "ns1", "0x12345").
		Return((&core.Verifier{
			Identity:  id.ID,
			Namespace: "ns1",
			VerifierRef: core.VerifierRef{
				Type:  core.VerifierTypeEthAddress,
				Value: "peer1",
			},
		}).Seal(), nil)
	mdi.On("GetIdentityByID", ctx, "ns1", id.ID).Return(nil, nil)

	_, err := im.cachedIdentityLookupByVerifierRef(ctx, "ns1", &core.VerifierRef{
		Type:  core.VerifierTypeEthAddress,
		Value: "0x12345",
	})
	assert.Regexp(t, "FF00116", err)

}

func TestCachedIdentityLookupMustExistCaching(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	id := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			DID:       "did:firefly:node/peer1",
			Namespace: "ns1",
			Name:      "peer1",
			Type:      core.IdentityTypeOrg,
		},
	}
	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetIdentityByDID", ctx, "ns1", "did:firefly:node/peer1").Return(id, nil).Once()

	v1, _, err := im.CachedIdentityLookupMustExist(ctx, "did:firefly:node/peer1")
	assert.NoError(t, err)
	assert.Equal(t, id, v1)

	v2, _, err := im.CachedIdentityLookupMustExist(ctx, "did:firefly:node/peer1")
	assert.NoError(t, err)
	assert.Equal(t, id, v2)
}

func TestCachedIdentityLookupMustExistUnknownResolver(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	_, retryable, err := im.CachedIdentityLookupMustExist(ctx, "did:random:anything")
	assert.Regexp(t, "FF10349", err)
	assert.False(t, retryable)

}

func TestCachedIdentityLookupMustExistGetIDFail(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetIdentityByDID", ctx, "ns1", "did:firefly:node/peer1").Return(nil, fmt.Errorf("pop"))

	_, retryable, err := im.CachedIdentityLookupMustExist(ctx, "did:firefly:node/peer1")
	assert.Regexp(t, "pop", err)
	assert.True(t, retryable)

}

func TestCachedIdentityLookupByVerifierByOldDIDFail(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	orgUUID := fftypes.NewUUID()
	did := core.FireFlyOrgDIDPrefix + orgUUID.String()

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetIdentityByDID", ctx, "ns1", did).Return(nil, nil)
	mdi.On("GetIdentityByID", ctx, "ns1", mock.MatchedBy(func(uuid *fftypes.UUID) bool {
		return uuid.Equals(orgUUID)
	})).Return(nil, fmt.Errorf("pop"))

	_, retryable, err := im.CachedIdentityLookupMustExist(ctx, did)
	assert.Regexp(t, "pop", err)
	assert.True(t, retryable)

}

func TestCachedIdentityLookupByIDCaching(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	id := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			DID:       "did:firefly:node/peer1",
			Namespace: "ns1",
			Name:      "peer1",
			Type:      core.IdentityTypeOrg,
		},
	}
	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetIdentityByID", ctx, "ns1", id.ID).Return(id, nil).Once()

	v1, err := im.CachedIdentityLookupByID(ctx, id.ID)
	assert.NoError(t, err)
	assert.Equal(t, id, v1)

	v2, err := im.CachedIdentityLookupByID(ctx, id.ID)
	assert.NoError(t, err)
	assert.Equal(t, id, v2)

	mdi.AssertExpectations(t)
}

func TestVerifyIdentityChainCustomOrgOrgOk(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	idRoot := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			DID:       "did:firefly:org/org1",
			Namespace: "ns1",
			Name:      "org1",
			Type:      core.IdentityTypeOrg,
		},
		Messages: core.IdentityMessages{
			Claim: fftypes.NewUUID(),
		},
	}
	idIntermediateOrg := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			Parent:    idRoot.ID,
			DID:       "did:firefly:org/org2",
			Namespace: "ns1",
			Name:      "org2",
			Type:      core.IdentityTypeOrg,
		},
		Messages: core.IdentityMessages{
			Claim: fftypes.NewUUID(),
		},
	}
	idIntermediateCustom := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			Parent:    idIntermediateOrg.ID,
			DID:       "did:firefly:ns/ns1/custom1",
			Namespace: "ns1",
			Name:      "custom1",
			Type:      core.IdentityTypeCustom,
		},
		Messages: core.IdentityMessages{
			Claim: fftypes.NewUUID(),
		},
	}
	idLeaf := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			Parent:    idIntermediateCustom.ID,
			DID:       "did:firefly:ns/ns1/custom2",
			Namespace: "ns1",
			Name:      "custom2",
			Type:      core.IdentityTypeCustom,
		},
	}
	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetIdentityByID", ctx, "ns1", idIntermediateOrg.ID).Return(idIntermediateOrg, nil).Once()
	mdi.On("GetIdentityByID", ctx, "ns1", idIntermediateCustom.ID).Return(idIntermediateCustom, nil).Once()
	mdi.On("GetIdentityByID", ctx, "ns1", idRoot.ID).Return(idRoot, nil).Once()

	immeidateParent, _, err := im.VerifyIdentityChain(ctx, idLeaf)
	assert.Equal(t, idIntermediateCustom, immeidateParent)
	assert.NoError(t, err)

	mdi.AssertExpectations(t)
}

func TestVerifyIdentityInvalid(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	id1 := &core.Identity{
		IdentityBase: core.IdentityBase{},
	}

	_, retryable, err := im.VerifyIdentityChain(ctx, id1)
	assert.Regexp(t, "FF00114", err)
	assert.False(t, retryable)

}

func TestVerifyIdentityChainLoop(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	idID1 := fftypes.NewUUID()
	idID2 := fftypes.NewUUID()
	id1 := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        idID1,
			Parent:    idID2,
			DID:       "did:firefly:org/org1",
			Namespace: "ns1",
			Name:      "org1",
			Type:      core.IdentityTypeOrg,
		},
		Messages: core.IdentityMessages{
			Claim: fftypes.NewUUID(),
		},
	}
	id2 := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        idID2,
			Parent:    idID1,
			DID:       "did:firefly:org/org2",
			Namespace: "ns1",
			Name:      "org2",
			Type:      core.IdentityTypeOrg,
		},
		Messages: core.IdentityMessages{
			Claim: fftypes.NewUUID(),
		},
	}

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetIdentityByID", ctx, "ns1", idID2).Return(id2, nil).Once()

	_, retryable, err := im.VerifyIdentityChain(ctx, id1)
	assert.Regexp(t, "FF10364", err)
	assert.False(t, retryable)

	mdi.AssertExpectations(t)
}

func TestVerifyIdentityChainBadParent(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	idID1 := fftypes.NewUUID()
	idID2 := fftypes.NewUUID()
	id1 := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        idID1,
			Parent:    idID2,
			DID:       "did:firefly:org/org1",
			Namespace: "ns1",
			Name:      "org1",
			Type:      core.IdentityTypeOrg,
		},
	}
	id2 := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        idID2,
			DID:       "did:firefly:org/org2",
			Namespace: "ns1",
			Name:      "org2",
			Type:      core.IdentityTypeOrg,
		},
	}

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetIdentityByID", ctx, "ns1", idID2).Return(id2, nil).Once()

	_, retryable, err := im.VerifyIdentityChain(ctx, id1)
	assert.Regexp(t, "FF10366", err)
	assert.False(t, retryable)

	mdi.AssertExpectations(t)
}

func TestVerifyIdentityChainErr(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	idID2 := fftypes.NewUUID()
	id1 := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			Parent:    idID2,
			DID:       "did:firefly:org/org1",
			Namespace: "ns1",
			Name:      "org1",
			Type:      core.IdentityTypeOrg,
		},
	}

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetIdentityByID", ctx, "ns1", idID2).Return(nil, fmt.Errorf("pop"))

	_, retryable, err := im.VerifyIdentityChain(ctx, id1)
	assert.Regexp(t, "pop", err)
	assert.True(t, retryable)

	mdi.AssertExpectations(t)
}

func TestVerifyIdentityChainNotFound(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	idID2 := fftypes.NewUUID()
	id1 := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			Parent:    idID2,
			DID:       "did:firefly:org/org1",
			Namespace: "ns1",
			Name:      "org1",
			Type:      core.IdentityTypeOrg,
		},
	}

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetIdentityByID", ctx, "ns1", idID2).Return(nil, nil)

	_, retryable, err := im.VerifyIdentityChain(ctx, id1)
	assert.Regexp(t, "FF10214", err)
	assert.False(t, retryable)

	mdi.AssertExpectations(t)
}

func TestVerifyIdentityChainInvalidParent(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	id1 := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			Parent:    nil,
			DID:       "did:firefly:ns/ns1/custom1",
			Namespace: "ns1",
			Name:      "custom1",
			Type:      core.IdentityTypeCustom,
		},
	}
	id2 := &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			Parent:    id1.ID,
			DID:       "did:firefly:org/org2",
			Namespace: "ns1",
			Name:      "org2",
			Type:      core.IdentityTypeOrg,
		},
	}

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetIdentityByID", ctx, "ns1", id1.ID).Return(id1, nil).Once()

	_, retryable, err := im.VerifyIdentityChain(ctx, id2)
	assert.Regexp(t, "FF10365", err)
	assert.False(t, retryable)

	mdi.AssertExpectations(t)
}

func TestValidateParentTypeCustomToNode(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	id1 := &core.Identity{
		IdentityBase: core.IdentityBase{
			Type: core.IdentityTypeNode,
		},
	}
	id2 := &core.Identity{
		IdentityBase: core.IdentityBase{
			Type: core.IdentityTypeCustom,
		},
	}

	err := im.validateParentType(ctx, id2, id1)
	assert.Regexp(t, "FF10365", err)

}

func TestValidateParentTypeInvalidType(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	id1 := &core.Identity{
		IdentityBase: core.IdentityBase{
			Type: core.IdentityTypeCustom,
		},
	}
	id2 := &core.Identity{
		IdentityBase: core.IdentityBase{
			Type: core.IdentityType("unknown"),
		},
	}

	err := im.validateParentType(ctx, id2, id1)
	assert.Regexp(t, "FF00126", err)

}

func TestCachedVerifierLookupCaching(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	verifier := (&core.Verifier{
		Namespace: "ns1",
		VerifierRef: core.VerifierRef{
			Value: "peer1",
			Type:  core.VerifierTypeFFDXPeerID,
		},
	}).Seal()
	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetVerifierByValue", ctx, verifier.Type, verifier.Namespace, verifier.Value).Return(verifier, nil).Once()

	v1, err := im.CachedVerifierLookup(ctx, core.VerifierTypeFFDXPeerID, "peer1")
	assert.NoError(t, err)
	assert.Equal(t, verifier, v1)

	v2, err := im.CachedVerifierLookup(ctx, core.VerifierTypeFFDXPeerID, "peer1")
	assert.NoError(t, err)
	assert.Equal(t, verifier, v2)

	mdi.AssertExpectations(t)
}

func TestCachedVerifierLookupError(t *testing.T) {

	ctx, im := newTestIdentityManager(t)

	mdi := im.database.(*databasemocks.Plugin)
	mdi.On("GetVerifierByValue", ctx, core.VerifierTypeFFDXPeerID, "ns1", "peer1").Return(nil, fmt.Errorf("pop"))

	_, err := im.CachedVerifierLookup(ctx, core.VerifierTypeFFDXPeerID, "peer1")
	assert.Regexp(t, "pop", err)

	mdi.AssertExpectations(t)
}

func TestResolveIdentitySignerOk(t *testing.T) {
	ctx, im := newTestIdentityManager(t)
	mdi := im.database.(*databasemocks.Plugin)

	msgID := fftypes.NewUUID()
	mdi.On("GetMessageByID", ctx, "ns1", msgID).Return(&core.Message{
		Header: core.MessageHeader{
			SignerRef: core.SignerRef{
				Key: "0x12345",
			},
		},
	}, nil)

	signerRef, err := im.ResolveIdentitySigner(ctx, &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			DID:       "did:firefly:org/org1",
			Namespace: "ns1",
			Name:      "org1",
			Type:      core.IdentityTypeOrg,
		},
		Messages: core.IdentityMessages{
			Claim: msgID,
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "0x12345", signerRef.Key)

	mdi.AssertExpectations(t)
}

func TestResolveIdentitySignerFail(t *testing.T) {
	ctx, im := newTestIdentityManager(t)
	mdi := im.database.(*databasemocks.Plugin)

	msgID := fftypes.NewUUID()
	mdi.On("GetMessageByID", ctx, "ns1", msgID).Return(nil, fmt.Errorf("pop"))

	_, err := im.ResolveIdentitySigner(ctx, &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			DID:       "did:firefly:org/org1",
			Namespace: "ns1",
			Name:      "org1",
			Type:      core.IdentityTypeOrg,
		},
		Messages: core.IdentityMessages{
			Claim: msgID,
		},
	})
	assert.Regexp(t, "pop", err)

	mdi.AssertExpectations(t)
}

func TestResolveIdentitySignerNotFound(t *testing.T) {
	ctx, im := newTestIdentityManager(t)
	mdi := im.database.(*databasemocks.Plugin)

	msgID := fftypes.NewUUID()
	mdi.On("GetMessageByID", ctx, "ns1", msgID).Return(nil, nil)

	_, err := im.ResolveIdentitySigner(ctx, &core.Identity{
		IdentityBase: core.IdentityBase{
			ID:        fftypes.NewUUID(),
			DID:       "did:firefly:org/org1",
			Namespace: "ns1",
			Name:      "org1",
			Type:      core.IdentityTypeOrg,
		},
		Messages: core.IdentityMessages{
			Claim: msgID,
		},
	})
	assert.Regexp(t, "FF10366", err)

	mdi.AssertExpectations(t)
}

func TestParseKeyNormalizationConfig(t *testing.T) {
	assert.Equal(t, KeyNormalizationBlockchainPlugin, ParseKeyNormalizationConfig("blockchain_Plugin"))
	assert.Equal(t, KeyNormalizationNone, ParseKeyNormalizationConfig("none"))
	assert.Equal(t, KeyNormalizationNone, ParseKeyNormalizationConfig(""))
}
