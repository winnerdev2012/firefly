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

package batch

import (
	"context"
	"database/sql/driver"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kaleido-io/firefly/internal/config"
	"github.com/kaleido-io/firefly/internal/database"
	"github.com/kaleido-io/firefly/internal/fftypes"
	"github.com/kaleido-io/firefly/internal/i18n"
	"github.com/kaleido-io/firefly/internal/log"
	"github.com/kaleido-io/firefly/internal/retry"
)

const (
	msgBatchOffsetName = "ff-msgbatch"
)

func NewBatchManager(ctx context.Context, database database.Plugin) (BatchManager, error) {
	if database == nil {
		return nil, i18n.NewError(ctx, i18n.MsgInitializationNilDepError)
	}
	readPageSize := config.GetUint(config.BatchManagerReadPageSize)
	bm := &batchManager{
		ctx:                        log.WithLogField(ctx, "role", "batchmgr"),
		database:                   database,
		readPageSize:               uint64(readPageSize),
		messagePollTimeout:         config.GetDuration(config.BatchManagerReadPollTimeout),
		startupOffsetRetryAttempts: config.GetInt(config.BatchManagerStartupAttempts),
		dispatchers:                make(map[fftypes.MessageType]*dispatcher),
		shoulderTap:                make(chan bool, 1),
		newMessages:                make(chan *uuid.UUID, readPageSize),
		sequencerClosed:            make(chan struct{}),
		retry: &retry.Retry{
			InitialDelay: config.GetDuration(config.BatchRetryInitDelay),
			MaximumDelay: config.GetDuration(config.BatchRetryMaxDelay),
			Factor:       config.GetFloat64(config.BatchRetryFactor),
		},
	}
	return bm, nil
}

type BatchManager interface {
	RegisterDispatcher(batchType fftypes.MessageType, handler DispatchHandler, batchOptions BatchOptions)
	NewMessages() chan<- *uuid.UUID
	Start() error
	Close()
	WaitStop()
}

type batchManager struct {
	ctx                        context.Context
	database                   database.Plugin
	dispatchers                map[fftypes.MessageType]*dispatcher
	shoulderTap                chan bool
	newMessages                chan *uuid.UUID
	sequencerClosed            chan struct{}
	retry                      *retry.Retry
	offset                     int64
	closed                     bool
	readPageSize               uint64
	messagePollTimeout         time.Duration
	startupOffsetRetryAttempts int
}

type DispatchHandler func(context.Context, *fftypes.Batch) error

type BatchOptions struct {
	BatchMaxSize   uint
	BatchTimeout   time.Duration
	DisposeTimeout time.Duration
}

type dispatcher struct {
	handler      DispatchHandler
	mux          sync.Mutex
	processors   map[string]*batchProcessor
	batchOptions BatchOptions
}

func (bm *batchManager) RegisterDispatcher(batchType fftypes.MessageType, handler DispatchHandler, batchOptions BatchOptions) {
	bm.dispatchers[batchType] = &dispatcher{
		handler:      handler,
		batchOptions: batchOptions,
		processors:   make(map[string]*batchProcessor),
	}
}

func (bm *batchManager) Start() error {
	if err := bm.restoreOffset(); err != nil {
		return err
	}
	go bm.newEventNotifications()
	go bm.messageSequencer()
	return nil
}

func (bm *batchManager) NewMessages() chan<- *uuid.UUID {
	return bm.newMessages
}

func (bm *batchManager) restoreOffset() error {
	offset, err := bm.database.GetOffset(bm.ctx, fftypes.OffsetTypeBatch, fftypes.SystemNamespace, msgBatchOffsetName)
	if err != nil {
		return err
	}
	if offset == nil {
		if err = bm.updateOffset(false, 0); err != nil {
			return err
		}
	} else {
		bm.offset = offset.Current
	}
	log.L(bm.ctx).Infof("Batch manager restored offset %d", bm.offset)
	return nil
}

func (bm *batchManager) removeProcessor(dispatcher *dispatcher, key string) {
	dispatcher.mux.Lock()
	delete(dispatcher.processors, key)
	dispatcher.mux.Unlock()
}

func (bm *batchManager) getProcessor(batchType fftypes.MessageType, namespace, author string) (*batchProcessor, error) {
	dispatcher, ok := bm.dispatchers[batchType]
	if !ok {
		return nil, i18n.NewError(bm.ctx, i18n.MsgUnregisteredBatchType, batchType)
	}
	dispatcher.mux.Lock()
	key := fmt.Sprintf("%s/%s", namespace, author)
	processor, ok := dispatcher.processors[key]
	if !ok {
		processor = newBatchProcessor(
			bm.ctx, // Background context, not the call context
			&batchProcessorConf{
				BatchOptions:       dispatcher.batchOptions,
				namespace:          namespace,
				author:             author,
				persitence:         bm.database,
				dispatch:           dispatcher.handler,
				processorQuiescing: func() { bm.removeProcessor(dispatcher, key) },
			},
			bm.retry,
		)
		dispatcher.processors[key] = processor
	}
	dispatcher.mux.Unlock()
	return processor, nil
}

func (bm *batchManager) Close() {
	if bm != nil && !bm.closed {
		for _, d := range bm.dispatchers {
			d.mux.Lock()
			for _, p := range d.processors {
				p.close()
			}
			d.mux.Unlock()
		}
		bm.closed = true
		close(bm.newMessages)
	}
	bm = nil
}

func (bm *batchManager) assembleMessageData(msg *fftypes.Message) (data []*fftypes.Data, err error) {
	// Load all the data - must all be present for us to send
	for _, dataRef := range msg.Data {
		if dataRef.ID == nil {
			continue
		}
		var d *fftypes.Data
		err = bm.retry.Do(bm.ctx, fmt.Sprintf("assemble %s data", dataRef.ID), func(attempt int) (retry bool, err error) {
			d, err = bm.database.GetDataById(bm.ctx, dataRef.ID)
			if err != nil {
				// continual retry for persistence error (distinct from not-found)
				return !bm.closed, err
			}
			return false, nil
		})
		if err != nil {
			return nil, err
		}
		if d == nil {
			return nil, i18n.NewError(bm.ctx, i18n.MsgDataNotFound, dataRef.ID)
		}
		data = append(data, d)
	}
	log.L(bm.ctx).Infof("Added broadcast message %s", msg.Header.ID)
	return data, nil
}

func (bm *batchManager) readPage() ([]*fftypes.Message, error) {
	var msgs []*fftypes.Message
	err := bm.retry.Do(bm.ctx, "retrieve messages", func(attempt int) (retry bool, err error) {
		fb := database.MessageQueryFactory.NewFilterLimit(bm.ctx, bm.readPageSize)
		msgs, err = bm.database.GetMessages(bm.ctx, fb.Gt("sequence", bm.offset).Sort("sequence").Limit(bm.readPageSize))
		if err != nil {
			return !bm.closed, err // Retry indefinitely, until closed (or context cancelled)
		}
		return false, nil
	})
	return msgs, err
}

func (bm *batchManager) messageSequencer() {
	l := log.L(bm.ctx)
	l.Debugf("Started batch assembly message sequencer")
	defer close(bm.sequencerClosed)

	dispatched := make(chan *batchDispatch, bm.readPageSize)

	for !bm.closed {
		// Read messages from the DB - in an error condition we retry until success, or a closed context
		msgs, err := bm.readPage()
		if err != nil {
			l.Debugf("Exiting: %s", err) // errors logged in readPage
			return
		}
		batchWasFull := false

		if len(msgs) > 0 {
			batchWasFull = (uint64(len(msgs)) == bm.readPageSize)
			var dispatchCount int
			for _, msg := range msgs {
				data, err := bm.assembleMessageData(msg)
				if err != nil {
					l.Errorf("Failed to retrieve message data for %s: %s", msg.Header.ID, err)
					continue
				}

				err = bm.dispatchMessage(dispatched, msg, data...)
				if err != nil {
					l.Errorf("Failed to dispatch message %s: %s", msg.Header.ID, err)
					continue
				}
				dispatchCount++
			}

			if dispatchCount > 0 {
				msgUpdates := make(map[uuid.UUID][]driver.Value)
				for i := 0; i < dispatchCount; i++ {
					dispatched := <-dispatched
					batchID := *dispatched.batchID
					l.Debugf("Dispatched message %s to batch %s", dispatched.msg.Header.ID, dispatched.batchID)
					msgUpdates[batchID] = append(msgUpdates[batchID], dispatched.msg.Header.ID)
				}
				if err = bm.updateMessages(msgUpdates); err != nil {
					l.Errorf("Closed while attempting to update messages: %s", err)
					l.Infof("Unflushed message updates: %+v", msgUpdates)
					break
				}
			}

			if !bm.closed {
				_ = bm.updateOffset(true, msgs[len(msgs)-1].Sequence)
			}
		}

		// Wait to be woken again
		if !bm.closed && !batchWasFull {
			bm.waitForShoulderTapOrPollTimeout()
		}
	}
}

// newEventNotifications just consumes new messags, logs them, then ensures there's a shoulderTap
// in the channel - without blocking. This is important as we must not block the notifier
func (bm *batchManager) newEventNotifications() {
	l := log.L(bm.ctx).WithField("role", "batch-newmessages")
	for {
		select {
		case m, ok := <-bm.newMessages:
			if !ok {
				l.Debugf("Exiting due to close")
				return
			}
			l.Debugf("Absorbing trigger for message %s", m)
		case <-bm.ctx.Done():
			l.Debugf("Exiting due to cancelled context")
			return
		}
		// Do not block sending to the shoulderTap - as it can only contain one
		select {
		case bm.shoulderTap <- true:
		default:
		}
	}
}

func (bm *batchManager) waitForShoulderTapOrPollTimeout() {
	l := log.L(bm.ctx)
	timeout := time.NewTimer(bm.messagePollTimeout)
	select {
	case <-timeout.C:
		l.Debugf("Woken after poll timeout")
	case <-bm.shoulderTap:
		l.Debugf("Woken for trigger for messages")
	case <-bm.ctx.Done():
		l.Debugf("Exiting due to cancelled context")
		bm.Close()
		return
	}
}

func (bm *batchManager) updateMessages(msgUpdates map[uuid.UUID][]driver.Value) (err error) {
	l := log.L(bm.ctx)
	return bm.retry.Do(bm.ctx, "update messages", func(attempt int) (retry bool, err error) {
		// Group the updates at the persistence layer
		err = bm.database.RunAsGroup(bm.ctx, func(ctx context.Context) error {
			// Group the updates by batch ID
			for batchID, msgs := range msgUpdates {
				f := database.MessageQueryFactory.NewFilter(ctx).In("id", msgs)
				u := database.MessageQueryFactory.NewUpdate(ctx).Set("batchid", &batchID)
				if err := bm.database.UpdateMessages(ctx, f, u); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			l.Errorf("Batch persist attempt %d failed: %s", attempt, err)
			return !bm.closed, err
		}
		return false, nil
	})
}

func (bm *batchManager) updateOffset(infiniteRetry bool, newOffset int64) (err error) {
	l := log.L(bm.ctx)
	return bm.retry.Do(bm.ctx, "update offset", func(attempt int) (retry bool, err error) {
		bm.offset = newOffset
		offset := &fftypes.Offset{
			Type:      fftypes.OffsetTypeBatch,
			Namespace: fftypes.SystemNamespace,
			Name:      msgBatchOffsetName,
			Current:   bm.offset,
		}
		err = bm.database.UpsertOffset(bm.ctx, offset)
		if err != nil {
			l.Errorf("Batch persist attempt %d failed: %s", attempt, err)
			stillRetrying := infiniteRetry || (attempt <= bm.startupOffsetRetryAttempts)
			return !bm.closed && stillRetrying, err
		}
		l.Infof("Batch manager committed offset %d", newOffset)
		return false, nil
	})
}

func (bm *batchManager) dispatchMessage(dispatched chan *batchDispatch, msg *fftypes.Message, data ...*fftypes.Data) error {
	l := log.L(bm.ctx)
	processor, err := bm.getProcessor(msg.Header.Type, msg.Header.Namespace, msg.Header.Author)
	if err != nil {
		return err
	}
	l.Debugf("Dispatching message %s to %s batch", msg.Header.ID, msg.Header.Type)
	work := &batchWork{
		msg:        msg,
		data:       data,
		dispatched: dispatched,
	}
	processor.newWork <- work
	return nil
}

func (bm *batchManager) WaitStop() {
	<-bm.sequencerClosed
	var processors []*batchProcessor
	for _, d := range bm.dispatchers {
		d.mux.Lock()
		for _, p := range d.processors {
			processors = append(processors, p)
		}
		d.mux.Unlock()
	}
	for _, p := range processors {
		p.waitClosed()
	}
}
