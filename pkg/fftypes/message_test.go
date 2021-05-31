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

package fftypes

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSealBareMessage(t *testing.T) {
	msg := Message{}
	err := msg.Seal(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, msg.Header.ID)
	assert.NotNil(t, msg.Header.DataHash)
	assert.NotNil(t, msg.Hash)
}

func TestSealNilDataID(t *testing.T) {
	msg := Message{
		Data: DataRefs{
			{},
		},
	}
	err := msg.Seal(context.Background())
	assert.Regexp(t, "FF10144.*0", err)
}

func TestVerifyNilDataHash(t *testing.T) {
	msg := Message{
		Data: DataRefs{
			{ID: NewUUID()},
		},
	}
	err := msg.Verify(context.Background())
	assert.Regexp(t, "FF10144.*0", err)
}

func TestSeaDupDataID(t *testing.T) {
	id1 := NewUUID()
	hash1 := NewRandB32()
	hash2 := NewRandB32()
	msg := Message{
		Data: DataRefs{
			{ID: id1, Hash: hash1},
			{ID: id1, Hash: hash2},
		},
	}
	err := msg.Seal(context.Background())
	assert.Regexp(t, "FF10145.*1", err)
}

func TestVerifylDupDataHash(t *testing.T) {
	id1 := NewUUID()
	id2 := NewUUID()
	hash1 := NewRandB32()
	msg := Message{
		Data: DataRefs{
			{ID: id1, Hash: hash1},
			{ID: id2, Hash: hash1},
		},
	}
	err := msg.Verify(context.Background())
	assert.Regexp(t, "FF10145.*1", err)
}

func TestVerifyNilHashes(t *testing.T) {
	msg := Message{}
	err := msg.Verify(context.Background())
	assert.Regexp(t, "FF10147", err)
}

func TestVerifyNilMisMatchedHashes(t *testing.T) {
	msg := Message{
		Header: MessageHeader{
			DataHash: NewRandB32(),
		},
		Hash: NewRandB32(),
	}
	err := msg.Verify(context.Background())
	assert.Regexp(t, "FF10146", err)
}

func TestSealKnownMessage(t *testing.T) {
	msgid := MustParseUUID("2cd37805-5f40-4e12-962e-67868cde3049")
	cid := MustParseUUID("39296b6e-91b9-4a61-b279-833c85b04d94")
	gid := MustParseUUID("5cd8afa6-f483-42f1-b11b-5a6f6421c81d")
	txid := MustParseUUID("87dbc29b-16e1-4578-bf24-0d3ac3b33ef1")
	data1 := MustParseUUID("e3a3b714-7e49-4c73-a4ea-87a50b19961a")
	data2 := MustParseUUID("cc66b23f-d340-4333-82d5-b63adc1c3c07")
	data3 := MustParseUUID("189c8185-2b92-481a-847a-e57595ab3541")
	var hash1, hash2, hash3 Bytes32
	hash1.UnmarshalText([]byte("3fcc7e07069e441f07c9f6b26f16fcb2dc896222d72888675082fd308440d9ae"))
	hash2.UnmarshalText([]byte("1d1462e02d7acee49a8448267c65067e0bec893c9a0c050b9835efa376fec046"))
	hash3.UnmarshalText([]byte("284b535da66aa0734af56c708426d756331baec3bce3079e508003bcf4738ee6"))
	msg := Message{
		Header: MessageHeader{
			ID:   msgid,
			CID:  cid,
			Type: MessageTypePrivate,
			TX: TransactionRef{
				Type: TransactionTypeBatchPin,
				ID:   txid,
			},
			Author:    "0x12345",
			Namespace: "ns1",
			Topic:     "topic1",
			Context:   "context1",
			Created:   UnixTime(1620104103123456789),
			Group:     gid,
		},
		Data: DataRefs{
			{ID: data1, Hash: &hash1},
			{Hash: &hash2, ID: data2}, // Demonstrating we hash in order id,hash order, even if Go definition is re-ordered
			{ID: data3, Hash: &hash3},
		},
	}
	err := msg.Seal(context.Background())
	assert.NoError(t, err)

	// Data IDs are hashed the the sequence of the message - that is preserved as the message moves through
	dataHashData, _ := json.Marshal(&msg.Data)
	var dataHash Bytes32 = sha256.Sum256([]byte(dataHashData))
	assert.Equal(t, `[{"id":"e3a3b714-7e49-4c73-a4ea-87a50b19961a","hash":"3fcc7e07069e441f07c9f6b26f16fcb2dc896222d72888675082fd308440d9ae"},{"id":"cc66b23f-d340-4333-82d5-b63adc1c3c07","hash":"1d1462e02d7acee49a8448267c65067e0bec893c9a0c050b9835efa376fec046"},{"id":"189c8185-2b92-481a-847a-e57595ab3541","hash":"284b535da66aa0734af56c708426d756331baec3bce3079e508003bcf4738ee6"}]`, string(dataHashData))
	assert.Equal(t, `2468d5c26cc85968acaf8b96d09476453916ea4eab41632a31d09efc7ab297d2`, dataHash.String())
	assert.Equal(t, dataHash, *msg.Header.DataHash)

	// Header contains the data hash, and is hashed into the message hash
	actualHeader, _ := json.Marshal(&msg.Header)
	expectedHeader := `{"id":"2cd37805-5f40-4e12-962e-67868cde3049","cid":"39296b6e-91b9-4a61-b279-833c85b04d94","type":"private","tx":{"type":"BatchPin","id":"87dbc29b-16e1-4578-bf24-0d3ac3b33ef1"},"author":"0x12345","created":"2021-05-04T04:55:03.123456789Z","namespace":"ns1","topic":"topic1","context":"context1","group":"5cd8afa6-f483-42f1-b11b-5a6f6421c81d","datahash":"2468d5c26cc85968acaf8b96d09476453916ea4eab41632a31d09efc7ab297d2"}`
	var msgHash Bytes32 = sha256.Sum256([]byte(expectedHeader))
	assert.Equal(t, expectedHeader, string(actualHeader))
	assert.Equal(t, `0375fc21a32d4b46338c26b394f72752fee0a85a40851bbb8d73479e7bfaea9f`, msgHash.String())
	assert.Equal(t, msgHash, *msg.Hash)

	// Verify also returns good
	err = msg.Verify(context.Background())
	assert.NoError(t, err)
}
