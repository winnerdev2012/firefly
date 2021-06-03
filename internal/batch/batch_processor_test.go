// Copyright © 2021 Kaleido, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in comdiliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or imdilied.
// See the License for the specific language governing permissions and
// limitations under the License.

package batch

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/kaleido-io/firefly/internal/log"
	"github.com/kaleido-io/firefly/internal/retry"
	"github.com/kaleido-io/firefly/mocks/databasemocks"
	"github.com/kaleido-io/firefly/pkg/fftypes"
	"github.com/likexian/gokit/assert"
	"github.com/stretchr/testify/mock"
)

func newTestBatchProcessor(dispatch DispatchHandler) (*databasemocks.Plugin, *batchProcessor) {
	mdi := &databasemocks.Plugin{}
	bp := newBatchProcessor(context.Background(), mdi, &batchProcessorConf{
		namespace:          "ns1",
		author:             "0x12345",
		dispatch:           dispatch,
		processorQuiescing: func() {},
		Options: Options{
			BatchMaxSize:   10,
			BatchTimeout:   10 * time.Millisecond,
			DisposeTimeout: 20 * time.Millisecond,
		},
	}, &retry.Retry{
		InitialDelay: 1 * time.Microsecond,
		MaximumDelay: 1 * time.Microsecond,
	})
	return mdi, bp
}

func mockRunAsGroupPassthrough(mdi *databasemocks.Plugin) {
	rag := mdi.On("RunAsGroup", mock.Anything, mock.Anything)
	rag.RunFn = func(a mock.Arguments) {
		fn := a[1].(func(context.Context) error)
		rag.ReturnArguments = mock.Arguments{fn(a[0].(context.Context))}
	}
}

func TestUnfilledBatch(t *testing.T) {
	log.SetLevel("debug")

	wg := sync.WaitGroup{}
	wg.Add(2)

	dispatched := []*fftypes.Batch{}
	mdi, bp := newTestBatchProcessor(func(c context.Context, b *fftypes.Batch) error {
		dispatched = append(dispatched, b)
		wg.Done()
		return nil
	})
	mockRunAsGroupPassthrough(mdi)
	mdi.On("UpdateMessages", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mdi.On("UpsertBatch", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mdi.On("UpdateBatch", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Generate the work the work
	work := make([]*batchWork, 5)
	for i := 0; i < 5; i++ {
		msgid := fftypes.NewUUID()
		work[i] = &batchWork{
			msg:        &fftypes.Message{Header: fftypes.MessageHeader{ID: msgid}},
			dispatched: make(chan *batchDispatch),
		}
	}

	// Kick off a go routine to consume the confirmations
	go func() {
		for i := 0; i < 5; i++ {
			<-work[i].dispatched
		}
		wg.Done()
	}()

	// Dispatch the work
	for i := 0; i < 5; i++ {
		bp.newWork <- work[i]
	}

	// Wait for the confirmations, and the dispatch
	wg.Wait()

	// Check we got all the messages in a single batch
	assert.Equal(t, len(dispatched[0].Payload.Messages), 5)

	bp.close()
	bp.waitClosed()

}

func TestFilledBatchSlowPersistence(t *testing.T) {
	log.SetLevel("debug")

	wg := sync.WaitGroup{}
	wg.Add(2)

	dispatched := []*fftypes.Batch{}
	mdi, bp := newTestBatchProcessor(func(c context.Context, b *fftypes.Batch) error {
		dispatched = append(dispatched, b)
		wg.Done()
		return nil
	})
	bp.conf.BatchTimeout = 1 * time.Hour // Must fill the batch
	mockUpsert := mdi.On("UpsertBatch", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	mockUpsert.ReturnArguments = mock.Arguments{nil}
	unblockPersistence := make(chan time.Time)
	mockUpsert.WaitFor = unblockPersistence
	mockRunAsGroupPassthrough(mdi)
	mdi.On("UpdateMessages", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mdi.On("UpdateBatch", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Generate the work the work
	work := make([]*batchWork, 10)
	for i := 0; i < 10; i++ {
		msgid := fftypes.NewUUID()
		if i%2 == 0 {
			work[i] = &batchWork{
				msg:        &fftypes.Message{Header: fftypes.MessageHeader{ID: msgid}},
				dispatched: make(chan *batchDispatch),
			}
		} else {
			work[i] = &batchWork{
				data:       []*fftypes.Data{{ID: msgid}},
				dispatched: make(chan *batchDispatch),
			}
		}
	}

	// Kick off a go routine to consume the confirmations
	go func() {
		for i := 0; i < 10; i++ {
			<-work[i].dispatched
		}
		wg.Done()
	}()

	// Dispatch the work
	for i := 0; i < 10; i++ {
		bp.newWork <- work[i]
	}

	// Unblock the dispatch
	time.Sleep(10 * time.Millisecond)
	mockUpsert.WaitFor = nil
	unblockPersistence <- time.Now() // First call to write the first entry in the batch

	// Wait for comdiletion
	wg.Wait()

	// Check we got all the messages in a single batch
	assert.Equal(t, len(dispatched[0].Payload.Messages), 5)
	assert.Equal(t, len(dispatched[0].Payload.Data), 5)

	bp.close()
	bp.waitClosed()

}

func TestCloseToUnblockDispatch(t *testing.T) {
	_, bp := newTestBatchProcessor(func(c context.Context, b *fftypes.Batch) error {
		return fmt.Errorf("pop")
	})
	bp.close()
	bp.dispatchBatch(&fftypes.Batch{})
}

func TestCloseToUnblockUpsertBatch(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(1)

	mdi, bp := newTestBatchProcessor(func(c context.Context, b *fftypes.Batch) error {
		return nil
	})
	bp.retry.MaximumDelay = 1 * time.Microsecond
	bp.conf.BatchTimeout = 100 * time.Second
	mockRunAsGroupPassthrough(mdi)
	mdi.On("UpdateMessages", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mup := mdi.On("UpsertBatch", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("pop"))
	waitForCall := make(chan bool)
	mup.RunFn = func(a mock.Arguments) {
		waitForCall <- true
		<-waitForCall
	}

	// Generate the work the work
	msgid := fftypes.NewUUID()
	work := &batchWork{
		msg:        &fftypes.Message{Header: fftypes.MessageHeader{ID: msgid}},
		dispatched: make(chan *batchDispatch),
	}

	// Dispatch the work
	bp.newWork <- work

	// Ensure the mock has been run
	<-waitForCall
	close(waitForCall)

	// Close to unblock
	bp.close()
	bp.waitClosed()

}
