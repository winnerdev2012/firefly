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

package cache

import (
	"context"
	"testing"
	"time"

	"github.com/hyperledger/firefly-common/pkg/config"
	"github.com/hyperledger/firefly/internal/coreconfig"
	"github.com/stretchr/testify/assert"
)

func TestNewCacheCreationFail(t *testing.T) {
	coreconfig.Reset()
	ctx := context.Background()
	cacheManager := NewCacheManager(ctx)
	_, err := cacheManager.GetCache(NewCacheConfig(ctx, "inconsistent1.limit", "inconsistent2.ttl", ""))
	assert.Equal(t, "FF10424: could not initialize cache - 'inconsistent1.limit' and 'inconsistent2.ttl' do not have identical prefix, mismatching prefixes are: 'inconsistent1','inconsistent2'", err.Error())
	_, err = cacheManager.GetCache(NewCacheConfig(ctx, "", "", ""))
	assert.Equal(t, "FF10422: could not initialize cache - size limit config key is not provided", err.Error())
	_, err = cacheManager.GetCache(NewCacheConfig(ctx, "test.limit", "", ""))
	assert.Equal(t, "FF10423: could not initialize cache - ttl config key is not provided", err.Error())
}

func TestGetCacheReturnsSameCacheForSameConfig(t *testing.T) {
	coreconfig.Reset()
	ctx := context.Background()
	cacheManager := NewCacheManager(ctx)
	cache0, _ := cacheManager.GetCache(NewCacheConfig(ctx, "cache.batch.limit", "cache.batch.ttl", "testnamespace"))
	cache1, _ := cacheManager.GetCache(NewCacheConfig(ctx, "cache.batch.limit", "cache.batch.ttl", "testnamespace"))

	assert.Equal(t, cache0, cache1)
	assert.Equal(t, []string{"testnamespace::cache.batch"}, cacheManager.ListKeys())

	cache2, _ := cacheManager.GetCache(NewCacheConfig(ctx, "cache.batch.limit", "cache.batch.ttl", ""))
	assert.NotEqual(t, cache0, cache2)
	assert.Equal(t, 2, len(cacheManager.ListKeys()))
}

func TestTwoSeparateCacheWorksIndependently(t *testing.T) {
	coreconfig.Reset()
	ctx := context.Background()
	cacheManager := NewCacheManager(ctx)
	cache0, _ := cacheManager.GetCache(NewCacheConfig(ctx, "cache.batch.limit", "cache.batch.ttl", ""))
	cache1, _ := cacheManager.GetCache(NewCacheConfig(ctx, "cache.operations.limit", "cache.operations.ttl", ""))

	cache0.SetInt("int0", 100)
	assert.Equal(t, 100, cache0.GetInt("int0"))
	assert.Equal(t, 100, cache0.Get("int0").(int))
	assert.Equal(t, 0, cache1.GetInt("int0"))
	assert.Equal(t, nil, cache1.Get("int0"))

	cache1.SetString("string1", "val1")
	assert.Equal(t, "", cache0.GetString("string1"))
	assert.Equal(t, nil, cache0.Get("string1"))
	assert.Equal(t, "val1", cache1.GetString("string1"))
	assert.Equal(t, "val1", cache1.Get("string1").(string))
}

func TestReturnsDummyCacheWhenCacheDisabled(t *testing.T) {
	coreconfig.Reset()
	config.Set(coreconfig.CacheEnabled, false)
	ctx := context.Background()
	cacheManager := NewCacheManager(ctx)
	cache0, _ := cacheManager.GetCache(NewCacheConfig(ctx, "cache.batch.limit", "cache.batch.ttl", ""))
	cache0.SetInt("int0", 100)
	assert.Equal(t, nil, cache0.Get("int0"))
}

func TestUmmanagedCacheInstance(t *testing.T) {
	uc0 := NewUmanagedCache(context.Background(), 100, 5*time.Minute)
	uc1 := NewUmanagedCache(context.Background(), 100, 5*time.Minute)
	assert.NotEqual(t, uc0, uc1)
}
