// Copyright © 2022 Kaleido, Inc.
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

package migrate

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteList(t *testing.T) {
	value := map[interface{}]interface{}{
		"values": []interface{}{"test1", "test2"},
	}
	config := &ConfigItem{value: value, writer: os.Stdout}
	config.Get("values").Each().Delete()
	assert.Equal(t, 0, len(value))
}

func TestNoRename(t *testing.T) {
	value := map[interface{}]interface{}{
		"key1": "val1",
		"key2": "val2",
	}
	config := &ConfigItem{value: value, writer: os.Stdout}
	config.Get("key1").RenameTo("key2")
	assert.Equal(t, map[interface{}]interface{}{
		"key2": "val2",
	}, value)
}
