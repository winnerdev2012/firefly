// Copyright © 2021 Kaleido, Inc.
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

package cmd

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"testing"

	"github.com/hyperledger/firefly/mocks/apiservermocks"
	"github.com/hyperledger/firefly/mocks/namespacemocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const configDir = "../test/data/config"

func TestGetEngine(t *testing.T) {
	assert.NotNil(t, getRootManager())
}

func TestExecMissingConfig(t *testing.T) {
	_utManager = &namespacemocks.Manager{}
	defer func() { _utManager = nil }()
	viper.Reset()
	err := Execute()
	assert.Regexp(t, "Not Found", err)
}

func TestShowConfig(t *testing.T) {
	_utManager = &namespacemocks.Manager{}
	defer func() { _utManager = nil }()
	viper.Reset()
	rootCmd.SetArgs([]string{"showconf"})
	defer rootCmd.SetArgs([]string{})
	err := rootCmd.Execute()
	assert.NoError(t, err)
}

func TestExecEngineInitFail(t *testing.T) {
	o := &namespacemocks.Manager{}
	o.On("Init", mock.Anything, mock.Anything).Return(fmt.Errorf("splutter"))
	_utManager = o
	defer func() { _utManager = nil }()
	os.Chdir(configDir)
	err := Execute()
	assert.Regexp(t, "splutter", err)
}

func TestExecEngineStartFail(t *testing.T) {
	o := &namespacemocks.Manager{}
	o.On("Init", mock.Anything, mock.Anything).Return(nil)
	o.On("Start").Return(fmt.Errorf("bang"))
	_utManager = o
	defer func() { _utManager = nil }()
	os.Chdir(configDir)
	err := Execute()
	assert.Regexp(t, "bang", err)
}

func TestExecOkExitSIGINT(t *testing.T) {
	o := &namespacemocks.Manager{}
	o.On("Init", mock.Anything, mock.Anything).Return(nil)
	o.On("Start").Return(nil)
	o.On("WaitStop").Return()
	_utManager = o
	defer func() { _utManager = nil }()

	os.Chdir(configDir)
	go func() {
		sigs <- syscall.SIGINT
	}()
	err := Execute()
	assert.NoError(t, err)
}

func TestAPIServerError(t *testing.T) {
	o := &namespacemocks.Manager{}
	o.On("Init", mock.Anything, mock.Anything).Return(nil)
	o.On("Start").Return(nil)
	as := &apiservermocks.Server{}
	as.On("Serve", mock.Anything, o).Return(fmt.Errorf("pop"))

	errChan := make(chan error)
	go startFirefly(context.Background(), o, as, errChan)
	err := <-errChan
	assert.EqualError(t, err, "pop")
}
