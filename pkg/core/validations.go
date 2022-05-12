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

package core

import (
	"context"
	"regexp"
	"strings"

	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-common/pkg/i18n"
)

var (
	ffNameValidator      = regexp.MustCompile(`^[0-9a-zA-Z]([0-9a-zA-Z._-]{0,62}[0-9a-zA-Z])?$`)
	ffSafeCharsValidator = regexp.MustCompile(`^[0-9a-zA-Z._-]*$`)
)

type NameValidationOptions struct {
	noUUID            bool
	extraAllowedChars string
}

func ValidateSafeCharsOnly(ctx context.Context, str string, fieldName string) error {
	if !ffSafeCharsValidator.MatchString(str) {
		return i18n.NewError(ctx, i18n.MsgSafeCharsOnly, fieldName)
	}
	return nil
}

func ValidateFFNameField(ctx context.Context, str string, fieldName string) error {
	return ValidateFFNameFieldOptions(ctx, str, fieldName, NameValidationOptions{})
}

func ValidateFFNameFieldNoUUID(ctx context.Context, str string, fieldName string) error {
	return ValidateFFNameFieldOptions(ctx, str, fieldName, NameValidationOptions{noUUID: true})
}

func ValidateFFNameFieldOptions(ctx context.Context, str string, fieldName string, options NameValidationOptions) error {
	if options.noUUID {
		// Name must not be a UUID
		if _, err := fftypes.ParseUUID(ctx, str); err == nil {
			return i18n.NewError(ctx, i18n.MsgNoUUID, fieldName)
		}
	}
	if options.extraAllowedChars != "" {
		// Remove some additional characters from validation
		for _, c := range options.extraAllowedChars {
			str = strings.ReplaceAll(str, string(c), "_")
		}
	}
	if !ffNameValidator.MatchString(str) {
		return i18n.NewError(ctx, i18n.MsgInvalidName, fieldName, options.extraAllowedChars)
	}
	return nil
}

func ValidateLength(ctx context.Context, str string, fieldName string, max int) error {
	if len([]byte(str)) > max {
		return i18n.NewError(ctx, i18n.MsgFieldTooLong, fieldName, max)
	}
	return nil
}
