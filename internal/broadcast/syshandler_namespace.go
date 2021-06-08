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

package broadcast

import (
	"context"

	"github.com/kaleido-io/firefly/internal/log"
	"github.com/kaleido-io/firefly/pkg/fftypes"
)

func (bm *broadcastManager) handleNamespaceBroadcast(ctx context.Context, msg *fftypes.Message, data []*fftypes.Data) (valid bool, err error) {
	l := log.L(ctx)

	var ns fftypes.Namespace
	valid = bm.getSystemBroadcastPayload(ctx, msg, data, &ns)
	if !valid {
		return false, nil
	}
	if err := ns.Validate(ctx, true); err != nil {
		l.Warnf("Unable to process namespace broadcast %s - validate failed: %s", msg.Header.ID, err)
		return false, nil
	}

	existing, err := bm.database.GetNamespace(ctx, ns.Name)
	if err != nil {
		return false, err // We only return database errors
	}
	if existing != nil {
		if existing.Type != fftypes.NamespaceTypeLocal {
			l.Warnf("Unable to process namespace broadcast %s (name=%s) - duplicate of %v", msg.Header.ID, existing.Name, existing.ID)
			return false, nil
		}
		// Remove the local definition
		if err = bm.database.DeleteNamespace(ctx, existing.ID); err != nil {
			return false, err
		}
	}

	if err = bm.database.UpsertNamespace(ctx, &ns, false); err != nil {
		return false, err
	}

	event := fftypes.NewEvent(fftypes.EventTypeNamespaceConfirmed, ns.Name, ns.ID, msg.Header.Group)
	if err = bm.database.UpsertEvent(ctx, event, false); err != nil {
		return false, err
	}

	return true, nil
}
