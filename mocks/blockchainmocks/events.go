// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package blockchainmocks

import (
	blockchain "github.com/kaleido-io/firefly/internal/blockchain"
	fftypes "github.com/kaleido-io/firefly/internal/fftypes"

	mock "github.com/stretchr/testify/mock"
)

// Events is an autogenerated mock type for the Events type
type Events struct {
	mock.Mock
}

// SequencedBroadcastBatch provides a mock function with given fields: batch, author, protocolTxId, additionalInfo
func (_m *Events) SequencedBroadcastBatch(batch *blockchain.BroadcastBatch, author string, protocolTxId string, additionalInfo map[string]interface{}) {
	_m.Called(batch, author, protocolTxId, additionalInfo)
}

// TransactionUpdate provides a mock function with given fields: txTrackingID, txState, protocolTxId, errorMessage, additionalInfo
func (_m *Events) TransactionUpdate(txTrackingID string, txState fftypes.TransactionState, protocolTxId string, errorMessage string, additionalInfo map[string]interface{}) {
	_m.Called(txTrackingID, txState, protocolTxId, errorMessage, additionalInfo)
}
