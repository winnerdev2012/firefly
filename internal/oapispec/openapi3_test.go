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

package oapispec

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/kaleido-io/firefly/internal/config"
	"github.com/kaleido-io/firefly/internal/i18n"
	"github.com/kaleido-io/firefly/pkg/database"
	"github.com/kaleido-io/firefly/pkg/fftypes"
	"github.com/stretchr/testify/assert"
)

func TestOpenAPI3SwaggerGen(t *testing.T) {

	config.Reset()
	routes := []*Route{
		{
			Name:   "op1",
			Path:   "namespaces/{ns}/example1/{id}",
			Method: http.MethodPost,
			PathParams: []*PathParam{
				{Name: "ns", ExampleFromConf: config.NamespacesDefault, Description: i18n.MsgTBD},
				{Name: "id", Description: i18n.MsgTBD},
			},
			QueryParams:     nil,
			FilterFactory:   nil,
			Description:     i18n.MsgTBD,
			JSONInputValue:  func() interface{} { return &fftypes.Message{} },
			JSONInputMask:   []string{"id"},
			JSONOutputValue: func() interface{} { return &fftypes.Batch{} },
			JSONOutputCode:  http.StatusOK,
		},
		{
			Name:            "op2",
			Path:            "example2",
			Method:          http.MethodGet,
			PathParams:      nil,
			QueryParams:     nil,
			FilterFactory:   database.MessageQueryFactory,
			Description:     i18n.MsgTBD,
			JSONInputValue:  func() interface{} { return nil },
			JSONOutputValue: func() interface{} { return []*fftypes.Batch{} },
			JSONOutputCode:  http.StatusOK,
		},
		{
			Name:       "op3",
			Path:       "example2",
			Method:     http.MethodPut,
			PathParams: nil,
			QueryParams: []*QueryParam{
				{Name: "ns", ExampleFromConf: config.NamespacesDefault, Description: i18n.MsgTBD},
				{Name: "id", Description: i18n.MsgTBD},
				{Name: "myfield", Default: "val1", Description: i18n.MsgTBD},
			},
			FilterFactory:   nil,
			Description:     i18n.MsgTBD,
			JSONInputValue:  func() interface{} { return &fftypes.Data{} },
			JSONOutputValue: func() interface{} { return nil },
			JSONOutputCode:  http.StatusNoContent,
		},
		{
			Name:   "op4",
			Path:   "example2/{id}",
			Method: http.MethodDelete,
			PathParams: []*PathParam{
				{Name: "id", Description: i18n.MsgTBD},
			},
			QueryParams:     nil,
			FilterFactory:   nil,
			Description:     i18n.MsgTBD,
			JSONInputValue:  func() interface{} { return nil },
			JSONOutputValue: func() interface{} { return nil },
			JSONOutputCode:  http.StatusNoContent,
		},
	}

	doc := SwaggerGen(context.Background(), routes)
	err := doc.Validate(context.Background())
	assert.NoError(t, err)

	b, err := yaml.Marshal(doc)
	assert.NoError(t, err)
	fmt.Print(string(b))

}
