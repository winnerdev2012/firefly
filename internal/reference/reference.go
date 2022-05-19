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

package reference

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-common/pkg/i18n"
	"github.com/hyperledger/firefly/internal/coremsgs"
	"github.com/hyperledger/firefly/pkg/core"
)

type TypeReferenceDoc struct {
	Example           []byte
	Description       []byte
	FieldDescriptions []byte
	SubFieldTables    []byte
}

/*
 * This function generates a series of markdown pages to document FireFly types, and are
 * designed to be included in the docs. Each page is a []byte value in the map, and the
 * key is the file name of the page. To add additional pages, simply create an example
 * instance of the type you would like to document, then include that in the `types`
 * array which is passed to generateMarkdownPages(). Note: It is the responsibility of
 * some other caller function to actually write the bytes to disk.
 */
func GenerateObjectsReferenceMarkdown(ctx context.Context) (map[string][]byte, error) {

	newest := core.SubOptsFirstEventNewest
	fifty := uint16(50)

	types := []interface{}{

		&core.Event{
			ID:          fftypes.MustParseUUID("5f875824-b36b-4559-9791-a57a2e2b30dd"),
			Sequence:    168,
			Type:        core.EventTypeTransactionSubmitted,
			Namespace:   "ns1",
			Reference:   fftypes.MustParseUUID("0d12aa75-5ed8-48a7-8b54-45274c6edcb1"),
			Transaction: fftypes.MustParseUUID("0d12aa75-5ed8-48a7-8b54-45274c6edcb1"),
			Topic:       core.TransactionTypeBatchPin.String(),
			Created:     fftypes.UnixTime(1652664195),
		},

		&core.Subscription{
			SubscriptionRef: core.SubscriptionRef{
				ID:        fftypes.MustParseUUID("c38d69fd-442e-4d6f-b5a4-bab1411c7fe8"),
				Namespace: "ns1",
				Name:      "app1",
			},
			Transport: "websockets",
			Filter: core.SubscriptionFilter{
				Events: "^(message_.*|token_.*)$",
				Message: core.MessageFilter{
					Tag: "^(red|blue)$",
				},
			},
			Options: core.SubscriptionOptions{
				SubscriptionCoreOptions: core.SubscriptionCoreOptions{
					FirstEvent: &newest,
					ReadAhead:  &fifty,
				},
			},
			Created: fftypes.UnixTime(1652664195),
		},

		&core.ContractAPI{
			ID:        fftypes.MustParseUUID("0f12317b-85a0-4a77-a722-857ea2b0a5fa"),
			Name:      "my_contract_api",
			Namespace: "ns1",
			Interface: &core.FFIReference{
				ID: fftypes.MustParseUUID("c35d3449-4f24-4676-8e64-91c9e46f06c4"),
			},
			Location: fftypes.JSONAnyPtr(`{
				"address": "0x95a6c4895c7806499ba35f75069198f45e88fc69"
			}`),
			Message: fftypes.MustParseUUID("b09d9f77-7b16-4760-a8d7-0e3c319b2a16"),
			URLs: core.ContractURLs{
				OpenAPI: "http://127.0.0.1:5000/api/v1/namespaces/default/apis/my_contract_api/api/swagger.json",
				UI:      "http://127.0.0.1:5000/api/v1/namespaces/default/apis/my_contract_api/api",
			},
		},

		&core.FFI{
			ID:          fftypes.MustParseUUID("c35d3449-4f24-4676-8e64-91c9e46f06c4"),
			Namespace:   "ns1",
			Name:        "SimpleStorage",
			Description: "A simple example contract in Solidity",
			Version:     "v0.0.1",
			Message:     fftypes.MustParseUUID("e4ad2077-5714-416e-81f9-7964a6223b6f"),
			Methods: []*core.FFIMethod{
				{
					ID:          fftypes.MustParseUUID("8f3289dd-3a19-4a9f-aab3-cb05289b013c"),
					Interface:   fftypes.MustParseUUID("c35d3449-4f24-4676-8e64-91c9e46f06c4"),
					Name:        "get",
					Namespace:   "ns1",
					Pathname:    "get",
					Description: "Get the current value",
					Params:      core.FFIParams{},
					Returns: core.FFIParams{
						{
							Name: "output",
							Schema: fftypes.JSONAnyPtr(`{
								"type": "integer",
								"details": {
								  "type": "uint256"
								}
							}`),
						},
					},
				},
				{
					ID:          fftypes.MustParseUUID("fc6f54ee-2e3c-4e56-b17c-4a1a0ae7394b"),
					Interface:   fftypes.MustParseUUID("c35d3449-4f24-4676-8e64-91c9e46f06c4"),
					Name:        "set",
					Namespace:   "ns1",
					Pathname:    "set",
					Description: "Set the value",
					Params: core.FFIParams{
						{
							Name: "newValue",
							Schema: fftypes.JSONAnyPtr(`{
								"type": "integer",
								"details": {
								  "type": "uint256"
								}
							}`),
						},
					},
					Returns: core.FFIParams{},
				},
			},
			Events: []*core.FFIEvent{
				{
					ID:        fftypes.MustParseUUID("9f653f93-86f4-45bc-be75-d7f5888fbbc0"),
					Interface: fftypes.MustParseUUID("c35d3449-4f24-4676-8e64-91c9e46f06c4"),
					Namespace: "ns1",
					Pathname:  "Changed",
					Signature: "Changed(address,uint256)",
					FFIEventDefinition: core.FFIEventDefinition{
						Name:        "Changed",
						Description: "Emitted when the value changes",
						Params: core.FFIParams{
							{
								Name: "_from",
								Schema: fftypes.JSONAnyPtr(`{
									"type": "string",
									"details": {
									  "type": "address",
									  "indexed": true
									}
								}`),
							},
							{
								Name: "_value",
								Schema: fftypes.JSONAnyPtr(`{
									"type": "integer",
									"details": {
									  "type": "uint256"
									}
								}`),
							},
						},
					},
				},
			},
		},

		&core.TokenPool{
			ID:        fftypes.MustParseUUID("90ebefdf-4230-48a5-9d07-c59751545859"),
			Type:      core.TokenTypeFungible,
			Namespace: "ns1",
			Name:      "my_token",
			Standard:  "ERC-20",
			Locator:   "address=0x056df1c53c3c00b0e13d37543f46930b42f71db0&schema=ERC20WithData&type=fungible",
			Decimals:  18,
			Connector: "erc20_erc721",
			State:     core.TokenPoolStateConfirmed,
			Message:   fftypes.MustParseUUID("43923040-b1e5-4164-aa20-47636c7177ee"),
			Info: fftypes.JSONObject{
				"address": "0x056df1c53c3c00b0e13d37543f46930b42f71db0",
				"name":    "pool8197",
				"schema":  "ERC20WithData",
			},
			TX: core.TransactionRef{
				Type: core.TransactionTypeTokenPool,
				ID:   fftypes.MustParseUUID("a23ffc87-81a2-4cbc-97d6-f53d320c36cd"),
			},
			Created: fftypes.UnixTime(1652664195),
		},

		&core.TokenTransfer{
			Message: fftypes.MustParseUUID("855af8e7-2b02-4e05-ad7d-9ae0d4c409ba"),
			Pool:    fftypes.MustParseUUID("1244ecbe-5862-41c3-99ec-4666a18b9dd5"),
			From:    "0x98151D8AB3af082A5DC07746C220Fb6C95Bc4a50",
			To:      "0x7b746b92869De61649d148823808653430682C0d",
			Type:    core.TokenTransferTypeTransfer,
			Created: fftypes.UnixTime(1652664195),
		},

		&core.Message{
			Header: core.MessageHeader{
				ID:     fftypes.MustParseUUID("4ea27cce-a103-4187-b318-f7b20fd87bf3"),
				Type:   core.MessageTypePrivate,
				CID:    fftypes.MustParseUUID("00d20cba-76ed-431d-b9ff-f04b4cbee55c"),
				TxType: core.TransactionTypeBatchPin,
				SignerRef: core.SignerRef{
					Author: "did:firefly:org/acme",
					Key:    "0xD53B0294B6a596D404809b1d51D1b4B3d1aD4945",
				},
				Created:   fftypes.UnixTime(1652664190),
				Group:     fftypes.HashString("testgroup"),
				Namespace: "ns1",
				Topics:    core.NewFFStringArray("topic1"),
				Tag:       "blue_message",
				DataHash:  fftypes.HashString("testmsghash"),
			},
			Data: []*core.DataRef{
				{
					ID:   fftypes.MustParseUUID("fdf9f118-eb81-4086-a63d-b06715b3bb4e"),
					Hash: fftypes.HashString("refhash"),
				},
			},
			State:     core.MessageStateConfirmed,
			Confirmed: fftypes.UnixTime(1652664196),
		},

		&core.Data{
			ID:        fftypes.MustParseUUID("4f11e022-01f4-4c3f-909f-5226947d9ef0"),
			Validator: core.ValidatorTypeJSON,
			Namespace: "ns1",
			Created:   fftypes.UnixTime(1652664195),
			Hash:      fftypes.HashString("testdatahash"),
			Datatype: &core.DatatypeRef{
				Name:    "widget",
				Version: "v1.2.3",
			},
			Value: fftypes.JSONAnyPtr(`{
				"name": "filename.pdf",
				"a": "example",
				"b": { "c": 12345 }
			}`),
			Blob: &core.BlobRef{
				Hash: fftypes.HashString("testblobhash"),
				Size: 12345,
				Name: "filename.pdf",
			},
		},

		&core.Batch{
			BatchHeader: core.BatchHeader{
				ID:        fftypes.MustParseUUID("894bc0ea-0c2e-4ca4-bbca-b4c39a816bbb"),
				Type:      core.BatchTypePrivate,
				Namespace: "ns1",
				Node:      fftypes.MustParseUUID("5802ab80-fa71-4f52-9189-fb534de93756"),
				Group:     fftypes.HashString("examplegroup"),
				Created:   fftypes.UnixTime(1652664196),
				SignerRef: core.SignerRef{
					Author: "did:firefly:org/example",
					Key:    "0x0a989907dcd17272257f3ebcf72f4351df65a846",
				},
			},
			Hash: fftypes.HashString("examplebatchhash"),
			Payload: core.BatchPayload{
				TX: core.TransactionRef{
					Type: core.BatchTypePrivate,
					ID:   fftypes.MustParseUUID("04930D84-0227-4044-9D6D-82C2952A0108"),
				},
				Messages: []*core.Message{},
				Data:     core.DataArray{},
			},
		},
	}

	simpleTypes := []interface{}{
		fftypes.UUID{},
		fftypes.FFTime{},
		fftypes.FFBigInt{},
		fftypes.JSONAny(""),
		fftypes.JSONObject{},
	}

	return generateMarkdownPages(ctx, types, simpleTypes, filepath.Join("..", "..", "docs", "reference", "types"))
}

func getType(v interface{}) reflect.Type {
	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		return reflect.TypeOf(v).Elem()
	}
	return reflect.TypeOf(v)
}

func generateMarkdownPages(ctx context.Context, types []interface{}, simpleTypes []interface{}, outputPath string) (map[string][]byte, error) {
	markdownMap := make(map[string][]byte, len(types))
	rootPageNames := make([]string, len(types))
	for i, v := range types {
		rootPageNames[i] = strings.ToLower(getType(v).Name())
	}

	simpleTypesMarkdown, simpleTypesNames := generateSimpleTypesMarkdown(ctx, simpleTypes, outputPath)
	markdownMap["simpletypes"] = simpleTypesMarkdown

	for i, o := range types {
		pageTitle := getType(types[i]).Name()
		// Page index starts at 1. Simple types will be the first page. Everything else comes after that.
		pageHeader := generatePageHeader(pageTitle, i+2)
		b := bytes.NewBuffer([]byte(pageHeader))
		markdown, _, err := generateObjectReferenceMarkdown(ctx, true, o, reflect.TypeOf(o), rootPageNames, simpleTypesNames, []string{}, outputPath)
		if err != nil {
			return nil, err
		}
		b.Write(markdown)
		markdownMap[rootPageNames[i]] = b.Bytes()
	}
	return markdownMap, nil
}

func generateSimpleTypesMarkdown(ctx context.Context, simpleTypes []interface{}, outputPath string) ([]byte, []string) {
	simpleTypeNames := make([]string, len(simpleTypes))
	for i, v := range simpleTypes {
		simpleTypeNames[i] = strings.ToLower(getType(v).Name())
	}

	pageHeader := generatePageHeader("Simple Types", 1)

	b := bytes.NewBuffer([]byte(pageHeader))
	for _, simpleType := range simpleTypes {
		markdown, _, _ := generateObjectReferenceMarkdown(ctx, true, nil, reflect.TypeOf(simpleType), []string{}, simpleTypeNames, []string{}, outputPath)
		b.Write(markdown)
	}
	return b.Bytes(), simpleTypeNames
}

func generateObjectReferenceMarkdown(ctx context.Context, descRequired bool, example interface{}, t reflect.Type, rootPageNames, simpleTypeNames, generatedTableNames []string, outputPath string) ([]byte, []string, error) {
	typeReferenceDoc := TypeReferenceDoc{}

	if t.Kind() == reflect.Ptr {
		t = reflect.TypeOf(example).Elem()
	}
	// generatedTableNames is where we keep track of all the tables we've generated (recursively)
	// for creating hyperlinks within the markdown
	generatedTableNames = append(generatedTableNames, strings.ToLower(t.Name()))

	// If a detailed type_description.md file exists, include that in a Description section here
	filename, _ := filepath.Abs(filepath.Join(outputPath, "includes", fmt.Sprintf("%s_description.md", strings.ToLower(t.Name()))))
	_, err := os.Stat(filename)
	if err != nil {
		if descRequired {
			return nil, nil, i18n.NewError(ctx, coremsgs.MsgReferenceMarkdownMissing, filename)
		}
	} else {
		typeReferenceDoc.Description = []byte(fmt.Sprintf("{%% include_relative includes/%s_description.md %%}\n\n", strings.ToLower(t.Name())))
	}

	// Include an example JSON representation if we have one available
	if example != nil {
		exampleJSON, err := json.MarshalIndent(example, "", "    ")
		if err != nil {
			return nil, nil, err
		}
		typeReferenceDoc.Example = exampleJSON
	}

	// If the type is a struct, look into each field inside it
	if t.Kind() == reflect.Struct {
		typeReferenceDoc.FieldDescriptions, typeReferenceDoc.SubFieldTables, generatedTableNames = generateFieldDescriptionsForStruct(ctx, t, rootPageNames, simpleTypeNames, generatedTableNames, outputPath)
	}

	// buff is the main buffer where we will write the markdown for this page
	buff := bytes.NewBuffer([]byte{})
	buff.WriteString(fmt.Sprintf("## %s\n\n", t.Name()))

	// If we only have one section, we will not write H3 headers
	sectionCount := 0
	if typeReferenceDoc.Description != nil {
		sectionCount++
	}
	if typeReferenceDoc.Example != nil {
		sectionCount++
	}
	if typeReferenceDoc.FieldDescriptions != nil {
		sectionCount++
	}

	if typeReferenceDoc.Description != nil {
		buff.Write(typeReferenceDoc.Description)
	}
	if typeReferenceDoc.Example != nil && len(typeReferenceDoc.Example) > 0 {
		if sectionCount > 1 {
			buff.WriteString("### Example\n\n```json\n")
		}
		buff.Write(typeReferenceDoc.Example)
		buff.WriteString("\n```\n\n")
	}
	if typeReferenceDoc.FieldDescriptions != nil && len(typeReferenceDoc.FieldDescriptions) > 0 {
		if sectionCount > 1 {
			buff.WriteString("### Field Descriptions\n\n")
		}
		buff.Write(typeReferenceDoc.FieldDescriptions)
		buff.WriteString("\n")
	}

	if typeReferenceDoc.SubFieldTables != nil && len(typeReferenceDoc.SubFieldTables) > 0 {
		buff.Write(typeReferenceDoc.SubFieldTables)
	}

	return buff.Bytes(), generatedTableNames, nil
}

func generateEnumList(f reflect.StructField) string {
	enumName := f.Tag.Get("ffenum")
	enumOptions := core.FFEnumValues(enumName)
	buff := new(strings.Builder)
	buff.WriteString("`FFEnum`:")
	for _, v := range enumOptions {
		buff.WriteString(fmt.Sprintf("<br/>`\"%s\"`", v))
	}
	return buff.String()
}

func generateFieldDescriptionsForStruct(ctx context.Context, t reflect.Type, rootPageNames, simpleTypeNames, generatedTableNames []string, outputPath string) ([]byte, []byte, []string) {
	fieldDescriptionsBytes := []byte{}
	// subFieldBuff is where we write any additional tables for sub fields that may be on this struct
	subFieldBuff := bytes.NewBuffer([]byte{})
	if t.NumField() > 0 {
		// Write the table to a temporary buffer - we will throw it away if there are no
		// public JSON serializable fields on the struct
		tableRowCount := 0
		tableBuff := bytes.NewBuffer([]byte{})
		tableBuff.WriteString("| Field Name | Description | Type |\n")
		tableBuff.WriteString("|------------|-------------|------|\n")

		tableRowCount = writeStructFields(ctx, t, rootPageNames, simpleTypeNames, generatedTableNames, outputPath, subFieldBuff, tableBuff, tableRowCount)

		if tableRowCount > 1 {
			fieldDescriptionsBytes = tableBuff.Bytes()
		}
	}
	return fieldDescriptionsBytes, subFieldBuff.Bytes(), generatedTableNames
}

func writeStructFields(ctx context.Context, t reflect.Type, rootPageNames, simpleTypeNames, generatedTableNames []string, outputPath string, subFieldBuff, tableBuff *bytes.Buffer, tableRowCount int) int {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		ffstructTag := field.Tag.Get("ffstruct")
		ffexcludeTag := field.Tag.Get("ffexclude")

		// If this is a nested struct, we need to recurse into it
		if field.Anonymous {
			tableRowCount = writeStructFields(ctx, field.Type, rootPageNames, simpleTypeNames, generatedTableNames, outputPath, subFieldBuff, tableBuff, tableRowCount)
			continue
		}

		// If the field is specifically excluded, or doesn't have a json tag, skip it
		if ffexcludeTag != "" || jsonTag == "" || jsonTag == "-" {
			continue
		}

		jsonFieldName := strings.Split(jsonTag, ",")[0]
		messageKeyName := fmt.Sprintf("%s.%s", ffstructTag, jsonFieldName)
		description := i18n.Expand(ctx, i18n.MessageKey(messageKeyName))
		isArray := false

		fieldType := field.Type
		fireflyType := fieldType.Name()

		if fieldType.Kind() == reflect.Slice {
			fieldType = fieldType.Elem()
			fireflyType = fieldType.Name()
			isArray = true
		}

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
			fireflyType = fieldType.Name()
		}

		if isArray {
			fireflyType = fmt.Sprintf("%s[]", fireflyType)
		}

		fireflyType = fmt.Sprintf("`%s`", fireflyType)

		isStruct := fieldType.Kind() == reflect.Struct
		isEnum := strings.ToLower(fieldType.Name()) == "ffenum"

		fieldInRootPages := false
		fieldInSimpleTypes := false
		for _, rootPageName := range rootPageNames {
			if strings.ToLower(fieldType.Name()) == rootPageName {
				fieldInRootPages = true
				break
			}
		}
		for _, simpleTypeName := range simpleTypeNames {
			if strings.ToLower(fieldType.Name()) == simpleTypeName {
				fieldInSimpleTypes = true
				break
			}
		}

		link := ""
		switch {
		case isEnum:
			fireflyType = generateEnumList(field)
		case fieldInRootPages:
			link = fmt.Sprintf("%s#%s", strings.ToLower(fieldType.Name()), strings.ToLower(fieldType.Name()))
		case fieldInSimpleTypes:
			link = fmt.Sprintf("simpletypes#%s", strings.ToLower(fieldType.Name()))
		case isStruct:
			link = fmt.Sprintf("#%s", strings.ToLower(fieldType.Name()))
		}
		if link != "" {
			fireflyType = fmt.Sprintf("[%s](%s)", fireflyType, link)

			// Generate the table for the sub type
			tableAlreadyGenerated := false
			for _, tableName := range generatedTableNames {
				if strings.ToLower(fieldType.Name()) == tableName {
					tableAlreadyGenerated = true
					break
				}
			}
			if isStruct && !tableAlreadyGenerated && !fieldInRootPages && !fieldInSimpleTypes {
				subFieldMarkdown, newTableNames, _ := generateObjectReferenceMarkdown(ctx, false, nil, fieldType, rootPageNames, simpleTypeNames, generatedTableNames, outputPath)
				generatedTableNames = newTableNames
				subFieldBuff.Write(subFieldMarkdown)
				subFieldBuff.WriteString("\n")
			}
		}

		tableBuff.WriteString(fmt.Sprintf("| %s | %s | %s |\n", jsonFieldName, description, fireflyType))
		tableRowCount++
	}
	return tableRowCount
}

func generatePageHeader(pageTitle string, navOrder int) string {
	return fmt.Sprintf(`---
layout: default
title: %s
parent: Core Resources
grand_parent: pages.reference
nav_order: %v
---

# %s
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---
`, pageTitle, navOrder, pageTitle)
}
