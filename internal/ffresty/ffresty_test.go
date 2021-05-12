// Copyright © 2021 Kaleido, Inc.
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

package ffresty

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/kaleido-io/firefly/internal/config"
	"github.com/kaleido-io/firefly/internal/i18n"
	"github.com/stretchr/testify/assert"
)

func TestRequestOK(t *testing.T) {

	customClient := &http.Client{}

	conf := config.NewPluginConfig("http_unit_tests")
	AddHTTPConfig(conf)
	conf.Set(HTTPConfigURL, "http://localhost:12345")
	conf.Set(HTTPConfigHeaders, map[string]interface{}{
		"someheader": "headervalue",
	})
	conf.Set(HTTPConfigAuthUsername, "user")
	conf.Set(HTTPConfigAuthPassword, "pass")
	conf.Set(HTTPCustomClient, customClient)
	defer config.Reset()

	c := New(context.Background(), conf)
	httpmock.ActivateNonDefault(customClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost:12345/test",
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "headervalue", req.Header.Get("someheader"))
			assert.Equal(t, "Basic dXNlcjpwYXNz", req.Header.Get("Authorization"))
			return httpmock.NewStringResponder(200, `{"some": "data"}`)(req)
		})

	resp, err := c.R().Get("/test")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, `{"some": "data"}`, resp.String())

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestRequestRetry(t *testing.T) {

	ctx := context.Background()

	conf := config.NewPluginConfig("http_unit_tests")
	AddHTTPConfig(conf)
	conf.Set(HTTPConfigURL, "http://localhost:12345")
	conf.Set(HTTPConfigRetryWaitTimeMS, 1)
	defer config.Reset()

	c := New(ctx, conf)
	httpmock.ActivateNonDefault(c.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost:12345/test",
		httpmock.NewStringResponder(500, `{"message": "pop"}`))

	resp, err := c.R().Get("/test")
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode())
	assert.Equal(t, 6, httpmock.GetTotalCallCount())

	err = WrapRestErr(ctx, resp, err, i18n.MsgEthconnectRESTErr)
	assert.Error(t, err)

}

func TestLongResponse(t *testing.T) {

	ctx := context.Background()

	conf := config.NewPluginConfig("http_unit_tests")
	AddHTTPConfig(conf)
	conf.Set(HTTPConfigURL, "http://localhost:12345")
	conf.Set(HTTPConfigRetryEnabled, false)
	defer config.Reset()

	c := New(ctx, conf)
	httpmock.ActivateNonDefault(c.GetClient())
	defer httpmock.DeactivateAndReset()

	resText := strings.Builder{}
	for i := 0; i < 512; i++ {
		resText.WriteByte(byte('a' + (i % 26)))
	}
	httpmock.RegisterResponder("GET", "http://localhost:12345/test",
		httpmock.NewStringResponder(500, resText.String()))

	resp, err := c.R().Get("/test")
	err = WrapRestErr(ctx, resp, err, i18n.MsgEthconnectRESTErr)
	assert.Error(t, err)
}

func TestErrResponse(t *testing.T) {

	ctx := context.Background()

	conf := config.NewPluginConfig("http_unit_tests")
	AddHTTPConfig(conf)
	conf.Set(HTTPConfigURL, "http://localhost:12345")
	conf.Set(HTTPConfigRetryEnabled, false)
	defer config.Reset()

	c := New(ctx, conf)
	httpmock.ActivateNonDefault(c.GetClient())
	defer httpmock.DeactivateAndReset()

	resText := strings.Builder{}
	for i := 0; i < 512; i++ {
		resText.WriteByte(byte('a' + (i % 26)))
	}
	httpmock.RegisterResponder("GET", "http://localhost:12345/test",
		httpmock.NewErrorResponder(fmt.Errorf("pop")))

	resp, err := c.R().Get("/test")
	err = WrapRestErr(ctx, resp, err, i18n.MsgEthconnectRESTErr)
	assert.Error(t, err)
}
