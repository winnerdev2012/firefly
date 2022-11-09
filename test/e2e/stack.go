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

package e2e

import (
	"encoding/json"
	"os"
)

type Stack struct {
	Name                  string    `json:"name,omitempty"`
	ExposedBlockchainPort int       `json:"exposedBlockchainPort,omitempty"`
	BlockchainProvider    string    `json:"blockchainProvider"`
	TokenProviders        []string  `json:"tokenProviders"`
	Members               []*Member `json:"members,omitempty"`
	Database              string    `json:"database"`
}

type StackState struct {
	Accounts []interface{} `json:"accounts"`
}

type Member struct {
	ExposedFireflyPort   int         `json:"exposedFireflyPort,omitempty"`
	ExposedAdminPort     int         `json:"exposedFireflyAdminPort,omitempty"`
	FireflyHostname      string      `json:"fireflyHostname,omitempty"`
	Username             string      `json:"username,omitempty"`
	Password             string      `json:"password,omitempty"`
	UseHTTPS             bool        `json:"useHttps,omitempty"`
	ExposedConnectorPort int         `json:"exposedConnectorPort,omitempty"`
	OrgName              string      `json:"orgName,omitempty"`
	Account              interface{} `json:"account,omitempty"`
}

func GetMemberPort(filename string, n int) (int, error) {
	jsonBytes, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	var stack Stack
	err = json.Unmarshal(jsonBytes, &stack)
	if err != nil {
		return 0, err
	}

	return stack.Members[n].ExposedFireflyPort, nil
}

func GetMemberHostname(filename string, n int) (string, error) {
	jsonBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	var stack Stack
	err = json.Unmarshal(jsonBytes, &stack)
	if err != nil {
		return "", err
	}

	return stack.Members[n].FireflyHostname, nil
}

func ReadStackFile(filename string) (*Stack, error) {
	jsonBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var stack *Stack
	err = json.Unmarshal(jsonBytes, &stack)
	if err != nil {
		return nil, err
	}

	// Apply defaults, in case this stack.json is a local CLI environment
	for _, member := range stack.Members {
		if member.FireflyHostname == "" {
			member.FireflyHostname = "127.0.0.1"
		}

	}

	return stack, nil
}

func ReadStackStateFile(filename string) (*StackState, error) {
	jsonBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var stackState *StackState
	err = json.Unmarshal(jsonBytes, &stackState)
	if err != nil {
		return nil, err
	}

	return stackState, nil
}
