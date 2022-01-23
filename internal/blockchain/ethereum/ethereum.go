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

package ethereum

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hyperledger/firefly/internal/config"
	"github.com/hyperledger/firefly/internal/config/wsconfig"
	"github.com/hyperledger/firefly/internal/i18n"
	"github.com/hyperledger/firefly/internal/log"
	"github.com/hyperledger/firefly/internal/restclient"
	"github.com/hyperledger/firefly/pkg/blockchain"
	"github.com/hyperledger/firefly/pkg/fftypes"
	"github.com/hyperledger/firefly/pkg/wsclient"
)

const (
	broadcastBatchEventSignature = "BatchPin(address,uint256,string,bytes32,bytes32,string,bytes32[])"
)

type Ethereum struct {
	ctx          context.Context
	topic        string
	instancePath string
	prefixShort  string
	prefixLong   string
	capabilities *blockchain.Capabilities
	callbacks    blockchain.Callbacks
	client       *resty.Client
	streams      *streamManager
	initInfo     struct {
		stream *eventStream
		sub    *subscription
	}
	wsconn          wsclient.WSClient
	closed          chan struct{}
	addressResolver *addressResolver
}

type eventStreamWebsocket struct {
	Topic string `json:"topic"`
}

type asyncTXSubmission struct {
	ID string `json:"id"`
}

type queryOutput struct {
	Output string `json:"output"`
}

type ethBatchPinInput struct {
	Namespace  string   `json:"namespace"`
	UUIDs      string   `json:"uuids"`
	BatchHash  string   `json:"batchHash"`
	PayloadRef string   `json:"payloadRef"`
	Contexts   []string `json:"contexts"`
}

type ethWSCommandPayload struct {
	Type  string `json:"type"`
	Topic string `json:"topic,omitempty"`
}

type Location struct {
	Address string `json:"address"`
}

type paramDetails struct {
	Type    string
	Indexed bool
}

var batchPinEvent = "BatchPin"
var addressVerify = regexp.MustCompile("^[0-9a-f]{40}$")

func (e *Ethereum) Name() string {
	return "ethereum"
}

func (e *Ethereum) Init(ctx context.Context, prefix config.Prefix, callbacks blockchain.Callbacks) (err error) {

	ethconnectConf := prefix.SubPrefix(EthconnectConfigKey)
	addressResolverConf := prefix.SubPrefix(AddressResolverConfigKey)

	e.ctx = log.WithLogField(ctx, "proto", "ethereum")
	e.callbacks = callbacks

	if addressResolverConf.GetString(AddressResolverURLTemplate) != "" {
		if e.addressResolver, err = newAddressResolver(ctx, addressResolverConf); err != nil {
			return err
		}
	}

	if ethconnectConf.GetString(restclient.HTTPConfigURL) == "" {
		return i18n.NewError(ctx, i18n.MsgMissingPluginConfig, "url", "blockchain.ethconnect")
	}
	e.instancePath = ethconnectConf.GetString(EthconnectConfigInstancePath)
	if e.instancePath == "" {
		return i18n.NewError(ctx, i18n.MsgMissingPluginConfig, "instance", "blockchain.ethconnect")
	}
	e.topic = ethconnectConf.GetString(EthconnectConfigTopic)
	if e.topic == "" {
		return i18n.NewError(ctx, i18n.MsgMissingPluginConfig, "topic", "blockchain.ethconnect")
	}

	e.prefixShort = ethconnectConf.GetString(EthconnectPrefixShort)
	e.prefixLong = ethconnectConf.GetString(EthconnectPrefixLong)

	e.client = restclient.New(e.ctx, ethconnectConf)
	e.capabilities = &blockchain.Capabilities{
		GlobalSequencer: true,
	}

	wsConfig := wsconfig.GenerateConfigFromPrefix(ethconnectConf)

	if wsConfig.WSKeyPath == "" {
		wsConfig.WSKeyPath = "/ws"
	}

	e.wsconn, err = wsclient.New(ctx, wsConfig, e.afterConnect)
	if err != nil {
		return err
	}

	e.streams = &streamManager{
		client: e.client,
	}
	batchSize := ethconnectConf.GetUint(EthconnectConfigBatchSize)
	batchTimeout := uint(ethconnectConf.GetDuration(EthconnectConfigBatchTimeout).Milliseconds())
	if e.initInfo.stream, err = e.streams.ensureEventStream(e.ctx, e.topic, batchSize, batchTimeout); err != nil {
		return err
	}
	log.L(e.ctx).Infof("Event stream: %s", e.initInfo.stream.ID)
	if e.initInfo.sub, err = e.streams.ensureSubscription(e.ctx, e.instancePath, e.initInfo.stream.ID, batchPinEvent); err != nil {
		return err
	}

	e.closed = make(chan struct{})
	go e.eventLoop()

	return nil
}

func (e *Ethereum) Start() error {
	return e.wsconn.Connect()
}

func (e *Ethereum) Capabilities() *blockchain.Capabilities {
	return e.capabilities
}

func (e *Ethereum) afterConnect(ctx context.Context, w wsclient.WSClient) error {
	// Send a subscribe to our topic after each connect/reconnect
	b, _ := json.Marshal(&ethWSCommandPayload{
		Type:  "listen",
		Topic: e.topic,
	})
	err := w.Send(ctx, b)
	if err == nil {
		b, _ = json.Marshal(&ethWSCommandPayload{
			Type: "listenreplies",
		})
		err = w.Send(ctx, b)
	}
	return err
}

func ethHexFormatB32(b *fftypes.Bytes32) string {
	if b == nil {
		return "0x0000000000000000000000000000000000000000000000000000000000000000"
	}
	return "0x" + hex.EncodeToString(b[0:32])
}

func (e *Ethereum) handleBatchPinEvent(ctx context.Context, msgJSON fftypes.JSONObject) (err error) {
	sBlockNumber := msgJSON.GetString("blockNumber")
	sTransactionIndex := msgJSON.GetString("transactionIndex")
	sTransactionHash := msgJSON.GetString("transactionHash")
	dataJSON := msgJSON.GetObject("data")
	authorAddress := dataJSON.GetString("author")
	ns := dataJSON.GetString("namespace")
	sUUIDs := dataJSON.GetString("uuids")
	sBatchHash := dataJSON.GetString("batchHash")
	sPayloadRef := dataJSON.GetString("payloadRef")
	sContexts := dataJSON.GetStringArray("contexts")
	timestampStr := msgJSON.GetString("timestamp")
	timestamp, err := fftypes.ParseTimeString(timestampStr)
	if err != nil {
		log.L(ctx).Errorf("BatchPin event is not valid - missing timestamp: %+v", msgJSON)
		return nil // move on
	}

	if sBlockNumber == "" ||
		sTransactionIndex == "" ||
		sTransactionHash == "" ||
		authorAddress == "" ||
		sUUIDs == "" ||
		sBatchHash == "" {
		log.L(ctx).Errorf("BatchPin event is not valid - missing data: %+v", msgJSON)
		return nil // move on
	}

	authorAddress, err = e.ResolveSigningKey(ctx, authorAddress)
	if err != nil {
		log.L(ctx).Errorf("BatchPin event is not valid - bad from address (%s): %+v", err, msgJSON)
		return nil // move on
	}

	hexUUIDs, err := hex.DecodeString(strings.TrimPrefix(sUUIDs, "0x"))
	if err != nil || len(hexUUIDs) != 32 {
		log.L(ctx).Errorf("BatchPin event is not valid - bad uuids (%s): %+v", err, msgJSON)
		return nil // move on
	}
	var txnID fftypes.UUID
	copy(txnID[:], hexUUIDs[0:16])
	var batchID fftypes.UUID
	copy(batchID[:], hexUUIDs[16:32])

	var batchHash fftypes.Bytes32
	err = batchHash.UnmarshalText([]byte(sBatchHash))
	if err != nil {
		log.L(ctx).Errorf("BatchPin event is not valid - bad batchHash (%s): %+v", err, msgJSON)
		return nil // move on
	}

	contexts := make([]*fftypes.Bytes32, len(sContexts))
	for i, sHash := range sContexts {
		var hash fftypes.Bytes32
		err = hash.UnmarshalText([]byte(sHash))
		if err != nil {
			log.L(ctx).Errorf("BatchPin event is not valid - bad pin %d (%s): %+v", i, err, msgJSON)
			return nil // move on
		}
		contexts[i] = &hash
	}

	delete(msgJSON, "data")
	batch := &blockchain.BatchPin{
		Namespace:       ns,
		TransactionID:   &txnID,
		BatchID:         &batchID,
		BatchHash:       &batchHash,
		BatchPayloadRef: sPayloadRef,
		Contexts:        contexts,
		Event: blockchain.Event{
			Source:     e.Name(),
			Name:       "BatchPin",
			ProtocolID: sTransactionHash,
			Output:     dataJSON,
			Info:       msgJSON,
			Timestamp:  timestamp,
		},
	}

	// If there's an error dispatching the event, we must return the error and shutdown
	return e.callbacks.BatchPinComplete(batch, authorAddress)
}

func (e *Ethereum) handleContractEvent(ctx context.Context, msgJSON fftypes.JSONObject) (err error) {
	sTransactionHash := msgJSON.GetString("transactionHash")
	sub := msgJSON.GetString("subId")
	signature := msgJSON.GetString("signature")
	dataJSON := msgJSON.GetObject("data")
	name := strings.SplitN(signature, "(", 2)[0]
	timestampStr := msgJSON.GetString("timestamp")
	timestamp, err := fftypes.ParseTimeString(timestampStr)
	if err != nil {
		log.L(ctx).Errorf("Contract event is not valid - missing timestamp: %+v", msgJSON)
		return err // move on
	}
	delete(msgJSON, "data")

	event := &blockchain.ContractEvent{
		Subscription: sub,
		Event: blockchain.Event{
			Source:     e.Name(),
			Name:       name,
			ProtocolID: sTransactionHash,
			Output:     dataJSON,
			Info:       msgJSON,
			Timestamp:  timestamp,
		},
	}

	return e.callbacks.ContractEvent(event)
}

func (e *Ethereum) handleReceipt(ctx context.Context, reply fftypes.JSONObject) error {
	l := log.L(ctx)

	headers := reply.GetObject("headers")
	requestID := headers.GetString("requestId")
	replyType := headers.GetString("type")
	txHash := reply.GetString("transactionHash")
	message := reply.GetString("errorMessage")
	if requestID == "" || replyType == "" {
		l.Errorf("Reply cannot be processed - missing fields: %+v", reply)
		return nil // Swallow this and move on
	}
	operationID, err := fftypes.ParseUUID(ctx, requestID)
	if err != nil {
		l.Errorf("Reply cannot be processed - bad ID: %+v", reply)
		return nil // Swallow this and move on
	}
	updateType := fftypes.OpStatusSucceeded
	if replyType != "TransactionSuccess" {
		updateType = fftypes.OpStatusFailed
	}
	l.Infof("Ethconnect '%s' reply: request=%s tx=%s message=%s", replyType, requestID, txHash, message)
	return e.callbacks.BlockchainOpUpdate(operationID, updateType, message, reply)
}

func (e *Ethereum) handleMessageBatch(ctx context.Context, messages []interface{}) error {
	l := log.L(ctx)

	for i, msgI := range messages {
		msgMap, ok := msgI.(map[string]interface{})
		if !ok {
			l.Errorf("Message cannot be parsed as JSON: %+v", msgI)
			return nil // Swallow this and move on
		}
		msgJSON := fftypes.JSONObject(msgMap)

		l1 := l.WithField("ethmsgidx", i)
		ctx1 := log.WithLogger(ctx, l1)
		signature := msgJSON.GetString("signature")
		sub := msgJSON.GetString("subId")
		l1.Infof("Received '%s' message", signature)
		l1.Tracef("Message: %+v", msgJSON)

		if sub == e.initInfo.sub.ID {
			switch signature {
			case broadcastBatchEventSignature:
				if err := e.handleBatchPinEvent(ctx1, msgJSON); err != nil {
					return err
				}
			default:
				l.Infof("Ignoring event with unknown signature: %s", signature)
			}
		} else if err := e.handleContractEvent(ctx1, msgJSON); err != nil {
			return err
		}
	}

	return nil
}

func (e *Ethereum) eventLoop() {
	defer e.wsconn.Close()
	defer close(e.closed)
	l := log.L(e.ctx).WithField("role", "event-loop")
	ctx := log.WithLogger(e.ctx, l)
	ack, _ := json.Marshal(map[string]string{"type": "ack", "topic": e.topic})
	for {
		select {
		case <-ctx.Done():
			l.Debugf("Event loop exiting (context cancelled)")
			return
		case msgBytes, ok := <-e.wsconn.Receive():
			if !ok {
				l.Debugf("Event loop exiting (receive channel closed)")
				return
			}

			var msgParsed interface{}
			err := json.Unmarshal(msgBytes, &msgParsed)
			if err != nil {
				l.Errorf("Message cannot be parsed as JSON: %s\n%s", err, string(msgBytes))
				continue // Swallow this and move on
			}
			switch msgTyped := msgParsed.(type) {
			case []interface{}:
				err = e.handleMessageBatch(ctx, msgTyped)
				if err == nil {
					err = e.wsconn.Send(ctx, ack)
				}
			case map[string]interface{}:
				err = e.handleReceipt(ctx, fftypes.JSONObject(msgTyped))
			default:
				l.Errorf("Message unexpected: %+v", msgTyped)
				continue
			}

			// Send the ack - only fails if shutting down
			if err != nil {
				l.Errorf("Event loop exiting: %s", err)
				return
			}
		}
	}
}

func varlidateEthAddress(ctx context.Context, key string) (string, error) {
	keyLower := strings.ToLower(key)
	keyNoHexPrefix := strings.TrimPrefix(keyLower, "0x")
	if addressVerify.MatchString(keyNoHexPrefix) {
		return "0x" + keyNoHexPrefix, nil
	}
	return "", i18n.NewError(ctx, i18n.MsgInvalidEthAddress)
}

func (e *Ethereum) ResolveSigningKey(ctx context.Context, key string) (string, error) {
	resolved, err := varlidateEthAddress(ctx, key)
	if err != nil && e.addressResolver != nil {
		resolved, err := e.addressResolver.ResolveSigningKey(ctx, key)
		if err == nil {
			log.L(ctx).Infof("Key '%s' resolved to '%s'", key, resolved)
		}
		return resolved, err
	}
	return resolved, err
}

func (e *Ethereum) invokeContractMethod(ctx context.Context, contractPath, method, signingKey, requestID string, input interface{}) (*resty.Response, error) {
	return e.client.R().
		SetContext(ctx).
		SetQueryParam(e.prefixShort+"-from", signingKey).
		SetQueryParam(e.prefixShort+"-sync", "false").
		SetQueryParam(e.prefixShort+"-id", requestID).
		SetBody(input).
		Post(contractPath + "/" + method)
}

func (e *Ethereum) queryContractMethod(ctx context.Context, contractPath, method string, input interface{}) (*resty.Response, error) {
	return e.client.R().
		SetContext(ctx).
		SetQueryParam(e.prefixShort+"-call", "true").
		SetBody(input).
		Post(contractPath + "/" + method)
}

func (e *Ethereum) SubmitBatchPin(ctx context.Context, operationID *fftypes.UUID, ledgerID *fftypes.UUID, signingKey string, batch *blockchain.BatchPin) error {
	ethHashes := make([]string, len(batch.Contexts))
	for i, v := range batch.Contexts {
		ethHashes[i] = ethHexFormatB32(v)
	}
	var uuids fftypes.Bytes32
	copy(uuids[0:16], (*batch.TransactionID)[:])
	copy(uuids[16:32], (*batch.BatchID)[:])
	input := &ethBatchPinInput{
		Namespace:  batch.Namespace,
		UUIDs:      ethHexFormatB32(&uuids),
		BatchHash:  ethHexFormatB32(batch.BatchHash),
		PayloadRef: batch.BatchPayloadRef,
		Contexts:   ethHashes,
	}
	res, err := e.invokeContractMethod(ctx, e.instancePath, "pinBatch", signingKey, operationID.String(), input)
	if err != nil || !res.IsSuccess() {
		return restclient.WrapRestErr(ctx, res, err, i18n.MsgEthconnectRESTErr)
	}
	return nil
}

func (e *Ethereum) InvokeContract(ctx context.Context, operationID *fftypes.UUID, signingKey string, location *fftypes.JSONAny, method *fftypes.FFIMethod, input map[string]interface{}) (interface{}, error) {
	contractAddress, err := parseContractLocation(ctx, location)
	if err != nil {
		return nil, err
	}
	res, err := e.invokeContractMethod(ctx, fmt.Sprintf("contracts/%v", contractAddress.Address), method.Name, signingKey, operationID.String(), input)
	if err != nil || !res.IsSuccess() {
		return nil, restclient.WrapRestErr(ctx, res, err, i18n.MsgEthconnectRESTErr)
	}
	tx := &asyncTXSubmission{}
	if err = json.Unmarshal(res.Body(), tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (e *Ethereum) QueryContract(ctx context.Context, location *fftypes.JSONAny, method *fftypes.FFIMethod, input map[string]interface{}) (interface{}, error) {
	contractAddress, err := parseContractLocation(ctx, location)
	if err != nil {
		return nil, err
	}
	res, err := e.queryContractMethod(ctx, fmt.Sprintf("contracts/%v", contractAddress.Address), method.Name, input)
	if err != nil || !res.IsSuccess() {
		return nil, restclient.WrapRestErr(ctx, res, err, i18n.MsgEthconnectRESTErr)
	}
	output := &queryOutput{}
	if err = json.Unmarshal(res.Body(), output); err != nil {
		return nil, err
	}
	return output, nil
}

func (e *Ethereum) ValidateContractLocation(ctx context.Context, location *fftypes.JSONAny) (err error) {
	_, err = parseContractLocation(ctx, location)
	return
}

func parseContractLocation(ctx context.Context, location *fftypes.JSONAny) (*Location, error) {
	ethLocation := Location{}
	if err := json.Unmarshal(location.Bytes(), &ethLocation); err != nil {
		return nil, i18n.NewError(ctx, i18n.MsgContractLocationInvalid, err)
	}
	if ethLocation.Address == "" {
		return nil, i18n.NewError(ctx, i18n.MsgContractLocationInvalid, "'address' not set")
	}
	return &ethLocation, nil
}

func parseParamDetails(ctx context.Context, details *fftypes.JSONAny) (*paramDetails, error) {
	ethParam := paramDetails{}
	if err := json.Unmarshal(details.Bytes(), &ethParam); err != nil {
		return nil, i18n.NewError(ctx, i18n.MsgContractParamInvalid, err)
	}
	if ethParam.Type == "" {
		return nil, i18n.NewError(ctx, i18n.MsgContractParamInvalid, "'type' not set")
	}
	return &ethParam, nil
}

var intRegex, _ = regexp.Compile("^u?int([0-9]{1,3})$")

func (e *Ethereum) ValidateFFIParam(ctx context.Context, param *fftypes.FFIParam) error {
	paramDetails, err := parseParamDetails(ctx, param.Details)
	if err != nil {
		return err
	}
	return e.validateParamInternal(ctx, param, paramDetails)
}

func (e *Ethereum) validateParamInternal(ctx context.Context, param *fftypes.FFIParam, paramDetails *paramDetails) error {
	switch {
	case len(param.Components) > 0:
		// struct
		if strings.HasPrefix(paramDetails.Type, "struct ") {
			for _, childParam := range param.Components {
				if err := e.ValidateFFIParam(ctx, childParam); err != nil {
					return err
				}
			}
			return nil
		}
	case strings.HasPrefix(param.Type, "byte"):
		// byte (array)
		if param.Type == paramDetails.Type {
			return nil
		}
		if paramDetails.Type == "byte[]" || strings.HasPrefix(paramDetails.Type, "bytes") {
			return nil
		}
	case strings.HasSuffix(param.Type, "[]"):
		// array
		if strings.Count(param.Type, "[]") == strings.Count(paramDetails.Type, "[]") {
			param.Type = strings.TrimSuffix(param.Type, "[]")
			paramDetails.Type = strings.TrimSuffix(paramDetails.Type, "[]")
			return e.validateParamInternal(ctx, param, paramDetails)
		}
	case param.Type == "integer":
		// integer
		matches := intRegex.FindStringSubmatch(paramDetails.Type)
		if len(matches) == 2 {
			i, err := strconv.ParseInt(matches[1], 10, 0)
			if err == nil && i >= 8 && i <= 256 && i%8 == 0 {
				return nil
			}
		}
	case param.Type == "string":
		// string
		if paramDetails.Type == "string" || paramDetails.Type == "address" {
			return nil
		}
	case param.Type == "boolean":
		if paramDetails.Type == "bool" {
			return nil
		}
	}
	return i18n.NewError(ctx, i18n.MsgContractInternalType, param.Name, param.Type, paramDetails.Type)
}

func (e *Ethereum) AddSubscription(ctx context.Context, subscription *fftypes.ContractSubscriptionInput) error {
	location, err := parseContractLocation(ctx, subscription.Location)
	if err != nil {
		return err
	}
	result, err := e.streams.createSubscription(ctx, location, e.initInfo.stream.ID, subscription.Event.FFIEventDefinition)
	if err != nil {
		return err
	}
	subscription.ProtocolID = result.ID
	return nil
}

func (e *Ethereum) DeleteSubscription(ctx context.Context, subscription *fftypes.ContractSubscription) error {
	return e.streams.deleteSubscription(ctx, subscription.ProtocolID)
}
