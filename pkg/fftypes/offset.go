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

package fftypes

type OffsetType string

const (
	// OffsetTypeBatch is an offset stored by the batch manager on the messages table
	OffsetTypeBatch OffsetType = "batch"
	// OffsetTypeAggregator is an offset stored by the aggregator on the events table
	OffsetTypeAggregator OffsetType = "aggregator"
	// OffsetTypeSubscription is an offeset stored by a dispatcher on the events table
	OffsetTypeSubscription OffsetType = "subscription"
)

// Offset is a simple stored data structure that records a sequence position within another collection
type Offset struct {
	ID        *UUID      `json:"id"`
	Type      OffsetType `json:"type"`
	Namespace string     `json:"namespace"`
	Name      string     `json:"name"`
	Current   int64      `json:"current,omitempty"`
}
