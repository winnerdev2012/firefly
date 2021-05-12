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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/kaleido-io/firefly/internal/config"
	"github.com/kaleido-io/firefly/internal/ffresty"
	"github.com/kaleido-io/firefly/internal/wsserver"
	"github.com/stretchr/testify/assert"
)

var utConfPrefix = config.NewPluginConfig("ws_unit_tests")

func resetConf() {
	config.Reset()
	InitConfigPrefix(utConfPrefix)
}

func TestWSClientE2E(t *testing.T) {

	wsServer := wsserver.NewWebSocketServer(context.Background())
	svr := httptest.NewServer(wsServer.Handler())
	defer svr.Close()

	acks := make(chan bool)
	sender, receiver, _ := wsServer.GetChannels("topic1")
	go func() {
		sender <- map[string]string{"test": "message"}
		err := <-receiver
		assert.Nil(t, err)
		acks <- true
	}()

	afterConnect := func(ctx context.Context, w WSClient) error {
		// Send a listen on topic1 in the connect options
		b, _ := json.Marshal(map[string]string{"type": "listen", "topic": "topic1"})
		return w.Send(ctx, b)
	}

	resetConf()
	utConfPrefix.Set(ffresty.HTTPConfigURL, fmt.Sprintf("ws://%s", svr.Listener.Addr()))
	utConfPrefix.Set(WSConfigKeyPath, "/ws")

	wsClient, err := New(context.Background(), utConfPrefix, afterConnect)
	assert.NoError(t, err)

	err = wsClient.Connect()
	assert.NoError(t, err)

	// Receive the message sent by the server
	b := <-wsClient.Receive()
	var msg map[string]string
	err = json.Unmarshal(b, &msg)
	assert.NoError(t, err)
	assert.Equal(t, "message", msg["test"])

	// Ack it
	b, _ = json.Marshal(map[string]string{"type": "ack", "topic": "topic1"})
	err = wsClient.Send(context.Background(), b)
	assert.NoError(t, err)

	// Wait for server to process our ack
	<-acks

	// Close out
	wsServer.Close()
	wsClient.Close()

}

func TestWSClientBadURL(t *testing.T) {
	resetConf()
	utConfPrefix.Set(ffresty.HTTPConfigURL, ":::")

	_, err := New(context.Background(), utConfPrefix, nil)
	assert.Regexp(t, "FF10162", err.Error())
}

func TestHTTPToWSURLRemap(t *testing.T) {
	resetConf()
	utConfPrefix.Set(ffresty.HTTPConfigURL, "http://test:12345")
	utConfPrefix.Set(WSConfigKeyPath, "/websocket")

	url, err := buildWSUrl(context.Background(), utConfPrefix)
	assert.NoError(t, err)
	assert.Equal(t, "ws://test:12345/websocket", url)
}

func TestHTTPSToWSSURLRemap(t *testing.T) {
	resetConf()
	utConfPrefix.Set(ffresty.HTTPConfigURL, "https://test:12345")

	url, err := buildWSUrl(context.Background(), utConfPrefix)
	assert.NoError(t, err)
	assert.Equal(t, "wss://test:12345", url)
}

func TestWSFailStartupHttp500(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "custom value", r.Header.Get("Custom-Header"))
			assert.Equal(t, "Basic dXNlcjpwYXNz", r.Header.Get("Authorization"))
			rw.WriteHeader(500)
			rw.Write([]byte(`{"error": "pop"}`))
		},
	))
	defer svr.Close()

	resetConf()
	utConfPrefix.Set(ffresty.HTTPConfigURL, fmt.Sprintf("ws://%s", svr.Listener.Addr()))
	utConfPrefix.Set(ffresty.HTTPConfigHeaders, map[string]interface{}{
		"custom-header": "custom value",
	})
	utConfPrefix.Set(ffresty.HTTPConfigAuthUsername, "user")
	utConfPrefix.Set(ffresty.HTTPConfigAuthPassword, "pass")
	utConfPrefix.Set(ffresty.HTTPConfigRetryWaitTimeMS, 1)
	utConfPrefix.Set(WSConfigKeyInitialConnectAttempts, 1)

	w, _ := New(context.Background(), utConfPrefix, nil)
	err := w.Connect()
	assert.Regexp(t, "FF10161", err.Error())
}

func TestWSFailStartupConnect(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(500)
		},
	))
	svr.Close()

	resetConf()
	utConfPrefix.Set(ffresty.HTTPConfigURL, fmt.Sprintf("ws://%s", svr.Listener.Addr()))
	utConfPrefix.Set(ffresty.HTTPConfigRetryWaitTimeMS, 1)
	utConfPrefix.Set(WSConfigKeyInitialConnectAttempts, 1)

	w, _ := New(context.Background(), utConfPrefix, nil)
	err := w.Connect()
	assert.Regexp(t, "FF10161", err.Error())
}

func TestWSSendClosed(t *testing.T) {

	wsServer := wsserver.NewWebSocketServer(context.Background())
	svr := httptest.NewServer(wsServer.Handler())
	defer svr.Close()

	resetConf()
	utConfPrefix.Set(ffresty.HTTPConfigURL, fmt.Sprintf("ws://%s", svr.Listener.Addr()))

	w, err := New(context.Background(), utConfPrefix, nil)
	assert.NoError(t, err)
	w.Close()

	err = w.Send(context.Background(), []byte(`sent after close`))
	assert.Regexp(t, "FF10160", err.Error())
}

func TestWSSendCancelledContext(t *testing.T) {

	w := &wsClient{
		send:    make(chan []byte),
		closing: make(chan struct{}),
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := w.Send(ctx, []byte(`sent after close`))
	assert.Regexp(t, "FF10159", err.Error())
}

func TestWSConnectClosed(t *testing.T) {

	w := &wsClient{
		ctx:    context.Background(),
		closed: true,
	}

	err := w.connect(false)
	assert.Regexp(t, "FF10160", err.Error())
}

func TestWSReadLoopSendFailure(t *testing.T) {

	wsServer := wsserver.NewWebSocketServer(context.Background())
	svr := httptest.NewServer(wsServer.Handler())
	defer svr.Close()
	defer wsServer.Close()

	sender, _, _ := wsServer.GetChannels("topic1")
	go func() {
		sender <- map[string]string{"test": "message"}
	}()

	wsconn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s", svr.Listener.Addr()), nil)
	wsconn.WriteJSON(map[string]string{"type": "listen", "topic": "topic1"})
	assert.NoError(t, err)
	defer wsconn.Close()
	w := &wsClient{
		ctx:      context.Background(),
		closed:   true,
		sendDone: make(chan []byte, 1),
		wsconn:   wsconn,
	}

	// Close the sender channel
	close(w.sendDone)

	// Ensure the readLoop exits immediately
	w.readLoop()

}

func TestWSReconnect(t *testing.T) {

	wsServer := wsserver.NewWebSocketServer(context.Background())
	svr := httptest.NewServer(wsServer.Handler())
	defer svr.Close()
	defer wsServer.Close()

	wsconn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s", svr.Listener.Addr()), nil)
	assert.NoError(t, err)
	wsconn.Close()
	ctxCancelled, cancel := context.WithCancel(context.Background())
	cancel()
	w := &wsClient{
		ctx:     ctxCancelled,
		receive: make(chan []byte),
		send:    make(chan []byte),
		closing: make(chan struct{}),
		wsconn:  wsconn,
	}
	close(w.send) // will mean sender exits immediately

	w.receiveReconnectLoop()
}

func TestWSSendFail(t *testing.T) {

	wsServer := wsserver.NewWebSocketServer(context.Background())
	svr := httptest.NewServer(wsServer.Handler())
	defer svr.Close()
	defer wsServer.Close()

	wsconn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s", svr.Listener.Addr()), nil)
	assert.NoError(t, err)
	wsconn.Close()
	w := &wsClient{
		ctx:      context.Background(),
		receive:  make(chan []byte),
		send:     make(chan []byte, 1),
		closing:  make(chan struct{}),
		sendDone: make(chan []byte, 1),
		wsconn:   wsconn,
	}
	w.send <- []byte(`wakes sender`)
	w.sendLoop()
	<-w.sendDone
}
