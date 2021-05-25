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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscriptionOptionsDatabaseSerialization(t *testing.T) {

	firstEvent := SubOptsFirstEventNewest
	readAhead := uint16(50)
	sub1 := &Subscription{
		Options: SubscriptionOptions{
			FirstEvent: &firstEvent,
			ReadAhead:  &readAhead,
		},
	}

	// Verify it serializes as bytes to the database
	b1, err := sub1.Options.Value()
	assert.NoError(t, err)
	assert.Equal(t, `{"firstEvent":"newest","readAhead":50}`, string(b1.([]byte)))

	// Verify it restores ok
	sub2 := &Subscription{}
	err = sub2.Options.Scan(b1)
	assert.NoError(t, err)
	b2, err := sub1.Options.Value()
	assert.NoError(t, err)
	assert.Equal(t, string(b1.([]byte)), string(b2.([]byte)))

	// Verify it can also scan as a string
	err = sub2.Options.Scan(string(b1.([]byte)))
	assert.NoError(t, err)

	// Out of luck with anything else
	err = sub2.Options.Scan(false)
	assert.Regexp(t, "FF10125", err.Error())

}
