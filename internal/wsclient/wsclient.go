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

package wsclient

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/kaleido-io/firefly/internal/config"
	"github.com/kaleido-io/firefly/internal/ffresty"
	"github.com/kaleido-io/firefly/internal/i18n"
	"github.com/kaleido-io/firefly/internal/log"
	"github.com/kaleido-io/firefly/internal/retry"
)

type WSAuthConfig struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type WSClient interface {
	Connect() error
	Receive() <-chan []byte
	Send(ctx context.Context, message []byte) error
	Close()
}

type wsClient struct {
	ctx                  context.Context
	headers              http.Header
	url                  string
	initialRetryAttempts int
	wsdialer             *websocket.Dialer
	wsconn               *websocket.Conn
	retry                retry.Retry
	closed               bool
	receive              chan []byte
	send                 chan []byte
	sendDone             chan []byte
	closing              chan struct{}
	afterConnect         WSPostConnectHandler
}

// WSPostConnectHandler will be called after every connect/reconnect. Can send data over ws, but must not block listening for data on the ws.
type WSPostConnectHandler func(ctx context.Context, w WSClient) error

func New(ctx context.Context, prefix config.ConfigPrefix, afterConnect WSPostConnectHandler) (WSClient, error) {

	wsURL, err := buildWSUrl(ctx, prefix)
	if err != nil {
		return nil, err
	}

	w := &wsClient{
		ctx: ctx,
		url: wsURL,
		wsdialer: &websocket.Dialer{
			ReadBufferSize:  prefix.GetInt(WSConfigKeyReadBufferSizeKB) * 1024,
			WriteBufferSize: prefix.GetInt(WSConfigKeyWriteBufferSizeKB) * 1024,
		},
		retry: retry.Retry{
			InitialDelay: prefix.GetDuration(ffresty.HTTPConfigRetryWaitTime),
			MaximumDelay: prefix.GetDuration(ffresty.HTTPConfigRetryMaxWaitTime),
		},
		initialRetryAttempts: prefix.GetInt(WSConfigKeyInitialConnectAttempts),
		headers:              make(http.Header),
		receive:              make(chan []byte),
		send:                 make(chan []byte),
		closing:              make(chan struct{}),
		afterConnect:         afterConnect,
	}
	for k, v := range prefix.GetObject(ffresty.HTTPConfigHeaders) {
		if vs, ok := v.(string); ok {
			w.headers.Set(k, vs)
		}
	}
	authUsername := prefix.GetString(ffresty.HTTPConfigAuthUsername)
	authPassword := prefix.GetString(ffresty.HTTPConfigAuthPassword)
	if authUsername != "" && authPassword != "" {
		w.headers.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", authUsername, authPassword)))))
	}

	return w, nil
}

func (w *wsClient) Connect() error {

	if err := w.connect(true); err != nil {
		return err
	}

	go w.receiveReconnectLoop()

	return nil
}

func (w *wsClient) Close() {
	if !w.closed {
		w.closed = true
		close(w.closing)
		c := w.wsconn
		if c != nil {
			_ = c.Close()
		}
	}
}

// Receive returns
func (w *wsClient) Receive() <-chan []byte {
	return w.receive
}

func (w *wsClient) Send(ctx context.Context, message []byte) error {
	// Send
	select {
	case w.send <- message:
		return nil
	case <-ctx.Done():
		return i18n.NewError(ctx, i18n.MsgWSSendTimedOut)
	case <-w.closing:
		return i18n.NewError(ctx, i18n.MsgWSClosing)
	}
}

func buildWSUrl(ctx context.Context, prefix config.ConfigPrefix) (string, error) {
	urlString := prefix.GetString(ffresty.HTTPConfigURL)
	u, err := url.Parse(urlString)
	if err != nil {
		return "", i18n.WrapError(ctx, err, i18n.MsgInvalidURL, urlString)
	}
	wsPath := prefix.GetString(WSConfigKeyPath)
	if wsPath != "" {
		u.Path = wsPath
	}
	if u.Scheme == "http" {
		u.Scheme = "ws"
	}
	if u.Scheme == "https" {
		u.Scheme = "wss"
	}
	return u.String(), nil
}

func (w *wsClient) connect(initial bool) error {
	l := log.L(w.ctx)
	return w.retry.Do(w.ctx, func(attempt int) (retry bool, err error) {
		if w.closed {
			return false, i18n.NewError(w.ctx, i18n.MsgWSClosing)
		}
		var res *http.Response
		w.wsconn, res, err = w.wsdialer.Dial(w.url, w.headers)
		if err != nil {
			var b []byte
			var status = -1
			if res != nil {
				b, _ = ioutil.ReadAll(res.Body)
				status = res.StatusCode
			}
			l.Warnf("WS %s connect attempt %d failed [%d]: %s", w.url, attempt, status, string(b))
			return !initial || attempt > w.initialRetryAttempts, i18n.WrapError(w.ctx, err, i18n.MsgWSConnectFailed)
		}
		l.Infof("WS %s connected", w.url)
		return false, nil
	})
}

func (w *wsClient) readLoop() {
	l := log.L(w.ctx)
	for {
		mt, message, err := w.wsconn.ReadMessage()

		// Check there's not a pending send message we need to return
		// before returning any error (do not block)
		select {
		case <-w.sendDone:
			l.Debugf("WS %s closing reader after send error", w.url)
			return
		default:
		}

		// return any error
		if err != nil {
			l.Errorf("WS %s closed: %s", w.url, err)
			return
		}

		// Pass the message to the consumer
		l.Tracef("WS %s read (mt=%d): %s", w.url, mt, message)
		w.receive <- message
	}
}

func (w *wsClient) sendLoop() {
	l := log.L(w.ctx)
	defer close(w.sendDone)

	for {
		message, ok := <-w.send
		if !ok {
			l.Debugf("WS %s send loop exiting", w.url)
			return
		}

		if err := w.wsconn.WriteMessage(websocket.TextMessage, message); err != nil {
			l.Errorf("WS %s send failed: %s", w.url, err)
			// Keep the message for when we reconnect
			w.sendDone <- message
			return
		}

	}
}

func (w *wsClient) receiveReconnectLoop() {
	l := log.L(w.ctx)
	defer close(w.receive)
	for !w.closed {
		// Start the sender, letting it close without blocking sending a notifiation on the sendDone
		w.sendDone = make(chan []byte, 1)
		go w.sendLoop()

		// Call the reconnect processor
		var err error
		if w.afterConnect != nil {
			err = w.afterConnect(w.ctx, w)
		}

		if err == nil {
			// Synchronously invoke the reader, as it's important we react immediately to any error there.
			w.readLoop()

			// Ensure the connection is closed after the sender exits
			err = w.wsconn.Close()
			if err != nil {
				l.Debugf("WS %s close failed: %s", w.url, err)
			}
			<-w.sendDone
			w.sendDone = nil
			w.wsconn = nil
		}

		// Go into reconnect
		if !w.closed {
			err = w.connect(false)
			if err != nil {
				l.Debugf("WS %s exiting: %s", w.url, err)
				return
			}
		}
	}
}
