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

package fftypes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrganizationValidation(t *testing.T) {

	org := &Organization{
		Name: "!name",
	}
	assert.Regexp(t, "FF10131.*name", org.Validate(context.Background(), false))

	org = &Organization{
		Name:        "ok",
		Description: string(make([]byte, 4097)),
	}
	assert.Regexp(t, "FF10188.*description", org.Validate(context.Background(), false))

	org = &Organization{
		Name:        "ok",
		Description: "ok",
		Identity:    "ok",
	}
	assert.NoError(t, org.Validate(context.Background(), false))

	assert.Regexp(t, "FF10203", org.Validate(context.Background(), true))

	var def Definition = org
	org.Identity = `A B C D E F G H I J K L M N O P Q R S T U V W X Y Z $ ( ) + ! 0 1 2 3 4 5 6 7 8 9`
	assert.Equal(t, "ff_org_A_B_C_D_E_F_G_H_I_J_K_L_M_N_O_P_Q_R_S_T_U_V_W_X_Y_Z___________0_1", def.Topic())
	def.SetBroadcastMessage(NewUUID())
	assert.NotNil(t, org.Message)
}
