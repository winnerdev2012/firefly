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

package apispec

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"github.com/kaleido-io/firefly/internal/config"
	"github.com/kaleido-io/firefly/internal/i18n"
)

func getHost(ctx context.Context) string {
	proto := "https"
	if !config.GetBool(config.HttpTLSEnabled) {
		proto = "http"
	}
	return fmt.Sprintf("%s://%s:%s", proto, config.GetString(config.HttpAddress), config.GetString(config.HttpPort))
}

func SwaggerGen(ctx context.Context, routes []*Route) *openapi3.T {

	doc := &openapi3.T{
		OpenAPI: "3.0.2",
		Servers: openapi3.Servers{
			{URL: fmt.Sprintf("%s/api/v1", getHost(ctx))},
		},
		Info: &openapi3.Info{
			Title:       "Firefly",
			Version:     "1.0",
			Description: "Copyright © 2021 Kaleido, Inc.",
		},
	}
	for _, route := range routes {
		addRoute(ctx, doc, route)
	}
	return doc
}

func getPathItem(ctx context.Context, doc *openapi3.T, path string) *openapi3.PathItem {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if doc.Paths == nil {
		doc.Paths = openapi3.Paths{}
	}
	pi, ok := doc.Paths[path]
	if ok {
		return pi
	}
	pi = &openapi3.PathItem{}
	doc.Paths[path] = pi
	return pi
}

func addInput(ctx context.Context, input interface{}, mask []string, op *openapi3.Operation) {
	schemaRef, _, _ := openapi3gen.NewSchemaRefForValue(maskFields(input, mask))
	op.RequestBody = &openapi3.RequestBodyRef{
		Value: &openapi3.RequestBody{
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: schemaRef,
				},
			},
		},
	}
}

func addOutput(ctx context.Context, route *Route, output interface{}, op *openapi3.Operation) {
	schemaRef, _, _ := openapi3gen.NewSchemaRefForValue(output)
	s := i18n.Expand(ctx, i18n.MsgSuccessResponse)
	op.Responses[strconv.FormatInt(int64(route.JSONOutputCode), 10)] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &s,
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: schemaRef,
				},
			},
		},
	}
}

func addParam(ctx context.Context, op *openapi3.Operation, in, name, def, example string, description i18n.MessageKey, msgArgs ...interface{}) {
	required := false
	if in == "path" {
		required = true
	}
	var defValue interface{} = nil
	if def != "" {
		defValue = &def
	}
	var exampleValue interface{}
	if example != "" {
		exampleValue = example
	}
	op.Parameters = append(op.Parameters, &openapi3.ParameterRef{
		Value: &openapi3.Parameter{
			In:          in,
			Name:        name,
			Required:    required,
			Description: i18n.Expand(ctx, description, msgArgs...),
			Schema: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:    "string",
					Default: defValue,
					Example: exampleValue,
				},
			},
		},
	})
}

func addRoute(ctx context.Context, doc *openapi3.T, route *Route) {
	pi := getPathItem(ctx, doc, route.Path)
	op := &openapi3.Operation{
		Description: i18n.Expand(ctx, route.Description),
		OperationID: route.Name,
		Responses:   openapi3.NewResponses(),
	}
	input := route.JSONInputValue()
	if input != nil {
		addInput(ctx, input, route.JSONInputMask, op)
	}
	output := route.JSONOutputValue()
	if output != nil {
		addOutput(ctx, route, output, op)
	}
	for _, p := range route.PathParams {
		example := p.Example
		if p.ExampleFromConf != "" {
			example = config.GetString(p.ExampleFromConf)
		}
		addParam(ctx, op, "path", p.Name, p.Default, example, p.Description)
	}
	for _, q := range route.QueryParams {
		example := q.Example
		if q.ExampleFromConf != "" {
			example = config.GetString(q.ExampleFromConf)
		}
		addParam(ctx, op, "query", q.Name, q.Default, example, q.Description)
	}
	if route.FilterFactory != nil {
		for _, field := range route.FilterFactory.NewFilter(ctx, 0).Fields() {
			addParam(ctx, op, "query", field, "", "", i18n.MsgFilterParamDesc)
		}
		addParam(ctx, op, "query", "sort", "", "", i18n.MsgFilterSortDesc)
		addParam(ctx, op, "query", "descending", "", "", i18n.MsgFilterDescendingDesc)
		addParam(ctx, op, "query", "skip", "", "", i18n.MsgFilterSkipDesc)
		addParam(ctx, op, "query", "limit", "", config.GetString(config.APIDefaultFilterLimit), i18n.MsgFilterLimitDesc)
	}
	switch route.Method {
	case http.MethodGet:
		pi.Get = op
	case http.MethodPut:
		pi.Put = op
	case http.MethodPost:
		pi.Post = op
	case http.MethodDelete:
		pi.Delete = op
	}
}

func maskFieldsOnStruct(t reflect.Type, mask []string) reflect.Type {
	fieldCount := t.NumField()
	newFields := make([]reflect.StructField, fieldCount)
	for i := 0; i < fieldCount; i++ {
		field := t.FieldByIndex([]int{i})
		if field.Type.Kind() == reflect.Struct {
			field.Type = maskFieldsOnStruct(field.Type, mask)
		} else {
			for _, m := range mask {
				if strings.EqualFold(field.Name, m) {
					field.Tag = "`json:-`"
				}
			}
		}
		newFields[i] = field
	}
	return reflect.StructOf(newFields)
}

func maskFields(input interface{}, mask []string) interface{} {
	t := reflect.TypeOf(input)
	newStruct := maskFieldsOnStruct(t.Elem(), mask)
	i := reflect.New(newStruct).Interface()
	return i
}
