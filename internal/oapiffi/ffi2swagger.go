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

package oapiffi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hyperledger/firefly/internal/i18n"
	"github.com/hyperledger/firefly/internal/oapispec"
	"github.com/hyperledger/firefly/pkg/fftypes"
)

type FFISwaggerGen interface {
	Generate(ctx context.Context, baseURL string, api *fftypes.ContractAPI, ffi *fftypes.FFI) *openapi3.T
}

// ffiSwaggerGen generates OpenAPI3 (Swagger) definitions for FFIs
type ffiSwaggerGen struct {
}

func NewFFISwaggerGen() FFISwaggerGen {
	return &ffiSwaggerGen{}
}

func (og *ffiSwaggerGen) Generate(ctx context.Context, baseURL string, api *fftypes.ContractAPI, ffi *fftypes.FFI) (swagger *openapi3.T) {
	hasLocation := !api.Location.IsNil()

	routes := []*oapispec.Route{
		{
			Name:            "apiRoot",
			Path:            "/",
			Method:          http.MethodGet,
			JSONInputValue:  nil,
			JSONOutputValue: func() interface{} { return &fftypes.ContractAPIWithInterface{} },
			JSONOutputCodes: []int{http.StatusOK},
			QueryParams: []*oapispec.QueryParam{
				{Name: "fetchinterface", IsBool: true, Description: i18n.MsgTBD, Example: "true"},
			},
		},
	}
	for _, method := range ffi.Methods {
		routes = og.addMethod(routes, method, hasLocation)
	}

	return oapispec.SwaggerGen(ctx, routes, &oapispec.SwaggerGenConfig{
		Title:       ffi.Name,
		Version:     ffi.Version,
		Description: ffi.Description,
		BaseURL:     baseURL,
	})
}

func (og *ffiSwaggerGen) addMethod(routes []*oapispec.Route, method *fftypes.FFIMethod, hasLocation bool) []*oapispec.Route {
	routes = append(routes, &oapispec.Route{
		Name:             fmt.Sprintf("invoke_%s", method.Pathname),
		Path:             fmt.Sprintf("invoke/%s", method.Pathname), // must match a route defined in apiserver routes!
		Method:           http.MethodPost,
		JSONInputSchema:  func(ctx context.Context) string { return contractCallJSONSchema(&method.Params, hasLocation).String() },
		JSONOutputSchema: func(ctx context.Context) string { return ffiParamsJSONSchema(&method.Returns).String() },
		JSONOutputCodes:  []int{http.StatusOK},
	})
	routes = append(routes, &oapispec.Route{
		Name:             fmt.Sprintf("query_%s", method.Pathname),
		Path:             fmt.Sprintf("query/%s", method.Pathname), // must match a route defined in apiserver routes!
		Method:           http.MethodPost,
		JSONOutputSchema: func(ctx context.Context) string { return ffiParamsJSONSchema(&method.Returns).String() },
		JSONOutputCodes:  []int{http.StatusOK},
	})
	return routes
}

/**
 * Parse the FFI and build a corresponding JSON Schema to describe the request body for "invoke".
 * Returns the JSON Schema as an `fftypes.JSONObject`.
 */
func contractCallJSONSchema(params *fftypes.FFIParams, hasLocation bool) *fftypes.JSONObject {
	req := &fftypes.ContractCallRequest{
		Input: *ffiParamsJSONSchema(params),
	}
	if !hasLocation {
		req.Location = fftypes.JSONAnyPtr(`{}`)
	}
	return &fftypes.JSONObject{
		"type":       "object",
		"properties": req,
	}
}

func ffiParamsJSONSchema(params *fftypes.FFIParams) *fftypes.JSONObject {
	out := make(fftypes.JSONObject, len(*params))
	for _, param := range *params {
		out[param.Name] = ffiParamJSONSchema(param)
	}
	return &fftypes.JSONObject{
		"type":       "object",
		"properties": out,
	}
}

func ffiParamJSONSchema(param *fftypes.FFIParam) *fftypes.JSONObject {
	out := fftypes.JSONObject{}
	if err := json.Unmarshal(param.Schema.Bytes(), &out); err == nil {
		return &out
	}
	return nil
}
