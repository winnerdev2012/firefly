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

package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hyperledger/firefly-common/pkg/config"
	"github.com/hyperledger/firefly-common/pkg/ffapi"
	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-common/pkg/i18n"
	"github.com/hyperledger/firefly/internal/coreconfig"
	"github.com/hyperledger/firefly/internal/coremsgs"
	"github.com/hyperledger/firefly/pkg/core"
	"github.com/hyperledger/firefly/pkg/database"
)

type FFISwaggerGen interface {
	Generate(ctx context.Context, baseURL string, api *core.ContractAPI, ffi *fftypes.FFI) *openapi3.T
}

type ContractListenerInput struct {
	Name    string                        `ffstruct:"ContractListener" json:"name,omitempty"`
	Topic   string                        `ffstruct:"ContractListener" json:"topic,omitempty"`
	Options *core.ContractListenerOptions `ffstruct:"ContractListener" json:"options,omitempty"`
}

type ContractListenerInputWithLocation struct {
	ContractListenerInput
	Location *fftypes.JSONAny `ffstruct:"ContractListener" json:"location,omitempty"`
}

// ffiSwaggerGen generates OpenAPI3 (Swagger) definitions for FFIs
type ffiSwaggerGen struct {
}

func NewFFISwaggerGen() FFISwaggerGen {
	return &ffiSwaggerGen{}
}

func (og *ffiSwaggerGen) Generate(ctx context.Context, baseURL string, api *core.ContractAPI, ffi *fftypes.FFI) (swagger *openapi3.T) {
	hasLocation := !api.Location.IsNil()

	routes := []*ffapi.Route{
		{
			Name:            "interface",
			Path:            "interface", // must match a route defined in apiserver routes!
			Method:          http.MethodGet,
			JSONInputValue:  nil,
			JSONOutputValue: func() interface{} { return &fftypes.FFI{} },
			JSONOutputCodes: []int{http.StatusOK},
		},
	}
	for _, method := range ffi.Methods {
		routes = og.addMethod(routes, method, hasLocation)
	}
	for _, event := range ffi.Events {
		routes = og.addEvent(routes, event, hasLocation)
	}

	return ffapi.NewSwaggerGen(&ffapi.Options{
		Title:                 ffi.Name,
		Version:               ffi.Version,
		Description:           ffi.Description,
		BaseURL:               baseURL,
		DefaultRequestTimeout: config.GetDuration(coreconfig.APIRequestTimeout),
	}).Generate(ctx, routes)
}

func (og *ffiSwaggerGen) addMethod(routes []*ffapi.Route, method *fftypes.FFIMethod, hasLocation bool) []*ffapi.Route {
	ctx := context.Background()
	description := method.Description
	if method.Details != nil && len(method.Details) > 0 {
		additionalDetailsHeader := i18n.Expand(ctx, coremsgs.APISmartContractDetails)
		description = fmt.Sprintf("%s\n\n%s:\n\n%s", description, additionalDetailsHeader, buildDetailsTable(ctx, method.Details))
	}
	routes = append(routes, &ffapi.Route{
		Name:   fmt.Sprintf("invoke_%s", method.Pathname),
		Path:   fmt.Sprintf("invoke/%s", method.Pathname), // must match a route defined in apiserver routes!
		Method: http.MethodPost,
		JSONInputSchema: func(ctx context.Context, schemaGen ffapi.SchemaGenerator) (*openapi3.SchemaRef, error) {
			return contractJSONSchema(&method.Params, hasLocation)
		},
		JSONOutputSchema: func(ctx context.Context, schemaGen ffapi.SchemaGenerator) (*openapi3.SchemaRef, error) {
			return contractJSONSchema(&method.Returns, true)
		},
		JSONOutputCodes:          []int{http.StatusOK},
		PreTranslatedDescription: description,
	})
	routes = append(routes, &ffapi.Route{
		Name:   fmt.Sprintf("query_%s", method.Pathname),
		Path:   fmt.Sprintf("query/%s", method.Pathname), // must match a route defined in apiserver routes!
		Method: http.MethodPost,
		JSONInputSchema: func(ctx context.Context, schemaGen ffapi.SchemaGenerator) (*openapi3.SchemaRef, error) {
			return contractJSONSchema(&method.Params, hasLocation)
		},
		JSONOutputSchema: func(ctx context.Context, schemaGen ffapi.SchemaGenerator) (*openapi3.SchemaRef, error) {
			return contractJSONSchema(&method.Returns, true)
		},
		JSONOutputCodes:          []int{http.StatusOK},
		PreTranslatedDescription: description,
	})
	return routes
}

func (og *ffiSwaggerGen) addEvent(routes []*ffapi.Route, event *fftypes.FFIEvent, hasLocation bool) []*ffapi.Route {
	ctx := context.Background()
	description := event.Description
	if event.Details != nil && len(event.Details) > 0 {
		additionalDetailsHeader := i18n.Expand(ctx, coremsgs.APISmartContractDetails)
		description = fmt.Sprintf("%s\n\n%s:\n\n%s", description, additionalDetailsHeader, buildDetailsTable(ctx, event.Details))
	}
	routes = append(routes, &ffapi.Route{
		Name:   fmt.Sprintf("createlistener_%s", event.Pathname),
		Path:   fmt.Sprintf("listeners/%s", event.Pathname), // must match a route defined in apiserver routes!
		Method: http.MethodPost,
		JSONInputValue: func() interface{} {
			if hasLocation {
				return &ContractListenerInput{}
			}
			return &ContractListenerInputWithLocation{}
		},
		JSONOutputValue:          func() interface{} { return &core.ContractListener{} },
		JSONOutputCodes:          []int{http.StatusOK},
		PreTranslatedDescription: description,
	})
	routes = append(routes, &ffapi.Route{
		Name:            fmt.Sprintf("getlistener_%s", event.Pathname),
		Path:            fmt.Sprintf("listeners/%s", event.Pathname), // must match a route defined in apiserver routes!
		Method:          http.MethodGet,
		JSONInputValue:  nil,
		JSONOutputValue: func() interface{} { return []*core.ContractListener{} },
		JSONOutputCodes: []int{http.StatusOK},
		Extensions: &coreExtensions{
			FilterFactory: database.ContractListenerQueryFactory,
		},
	})
	return routes
}

/**
 * Parse the FFI and build a corresponding JSON Schema to describe the request body for "invoke".
 * Returns the JSON Schema as an `fftypes.JSONObject`.
 */
func contractJSONSchema(params *fftypes.FFIParams, hasLocation bool) (*openapi3.SchemaRef, error) {
	paramSchema := make(fftypes.JSONObject, len(*params))
	for _, param := range *params {
		paramSchema[param.Name] = param.Schema
	}
	inputSchema := fftypes.JSONObject{
		"type":       "object",
		"properties": paramSchema,
	}
	properties := fftypes.JSONObject{
		"input": inputSchema,
		"options": fftypes.JSONObject{
			"type": "object",
		},
	}
	if !hasLocation {
		properties["location"] = fftypes.JSONAnyPtr(`{}`)
	}
	schema := fftypes.JSONObject{
		"type":       "object",
		"properties": properties,
	}
	b, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	s := openapi3.NewSchema()
	err = s.UnmarshalJSON(b)
	if err != nil {
		return nil, err
	}
	return openapi3.NewSchemaRef("", s), nil
}

func buildDetailsTable(ctx context.Context, details map[string]interface{}) string {
	keyHeader := i18n.Expand(ctx, coremsgs.APISmartContractDetailsKey)
	valueHeader := i18n.Expand(ctx, coremsgs.APISmartContractDetailsKey)
	var s strings.Builder
	s.WriteString(fmt.Sprintf("| %s | %s |\n|-----|-------|\n", keyHeader, valueHeader))
	keys := make([]string, len(details))
	i := 0
	for key := range details {
		keys[i] = key
	}
	sort.Strings(keys)
	for _, key := range keys {
		s.WriteString(fmt.Sprintf("|%s|%s|\n", key, details[key]))
	}
	return s.String()
}
