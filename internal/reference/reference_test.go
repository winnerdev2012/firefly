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

//go:build !reference
// +build !reference

package reference

import (
	"context"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckGeneratedMarkdownPages(t *testing.T) {
	markdownMap, err := GenerateObjectsReferenceMarkdown()
	assert.NoError(t, err)
	assert.NotNil(t, markdownMap)

	for pageName, markdown := range markdownMap {
		b, err := os.ReadFile(filepath.Join("..", "..", "docs", "reference", "types", fmt.Sprintf("%s.md", pageName)))
		assert.NoError(t, err)
		expectedPageHash := sha1.New()
		expectedPageHash.Write(b)
		actualPageHash := sha1.New()
		actualPageHash.Write(markdown)
		assert.Equal(t, expectedPageHash.Sum(nil), actualPageHash.Sum(nil), "The type reference docs generated by the code did not match the docs files in git for page: '%s.md' Did you forget to run `make reference`?", pageName)
	}
}

func TestGenerateMarkdownPagesNonPointer(t *testing.T) {
	markdownMap, err := generateMarkdownPages(context.Background(), []interface{}{"foo"}, "")
	assert.NoError(t, err)
	assert.NotNil(t, markdownMap)
}

func TestGenerateMarkdownPagesBadJSON(t *testing.T) {
	badJSON := map[bool]bool{true: false}
	_, err := generateMarkdownPages(context.Background(), []interface{}{badJSON}, "")
	assert.Error(t, err)
}
