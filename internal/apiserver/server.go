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

package apiserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/kaleido-io/firefly/internal/config"
	"github.com/kaleido-io/firefly/internal/engine"
	"github.com/kaleido-io/firefly/internal/fftypes"
	"github.com/kaleido-io/firefly/internal/i18n"
	"github.com/kaleido-io/firefly/internal/log"
)

var ffcodeExtractor = regexp.MustCompile(`^(FF\d+):`)

// Serve is the main entry point for the API Server
func Serve(ctx context.Context, initEngine bool) error {
	e := engine.NewEngine()
	if initEngine {
		if err := e.Init(ctx); err != nil {
			return err
		}
	}
	defer e.Close()

	r := createMuxRouter(e)
	l, err := createListener(ctx)
	if err == nil {
		var s *http.Server
		s, err = createServer(ctx, r)
		if err == nil {
			err = serveHTTP(ctx, l, s)
		}
	}
	return err
}

func createListener(ctx context.Context) (net.Listener, error) {
	listenAddr := fmt.Sprintf("%s:%d", config.GetString(config.HttpAddress), config.GetUint(config.HttpPort))
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, i18n.WrapError(ctx, err, i18n.MsgAPIServerStartFailed, listenAddr)
	}
	log.L(ctx).Infof("Listening on HTTP %s", listener.Addr())
	return listener, err
}

func createServer(ctx context.Context, r *mux.Router) (srv *http.Server, err error) {

	// Support client auth
	clientAuth := tls.NoClientCert
	if config.GetBool(config.HttpTLSClientAuth) {
		clientAuth = tls.RequireAndVerifyClientCert
	}

	// Support custom CA file
	var rootCAs *x509.CertPool
	caFile := config.GetString(config.HttpTLSCAFile)
	if caFile != "" {
		rootCAs = x509.NewCertPool()
		var caBytes []byte
		caBytes, err = ioutil.ReadFile(caFile)
		if err == nil {
			ok := rootCAs.AppendCertsFromPEM(caBytes)
			if !ok {
				err = i18n.NewError(ctx, i18n.MsgInvalidCAFile)
			}
		}
	} else {
		rootCAs, err = x509.SystemCertPool()
	}

	if err != nil {
		return nil, i18n.WrapError(ctx, err, i18n.MsgTLSConfigFailed)
	}

	srv = &http.Server{
		Handler:      r,
		WriteTimeout: time.Duration(config.GetUint(config.HttpWriteTimeout)) * time.Millisecond,
		ReadTimeout:  time.Duration(config.GetUint(config.HttpReadTimeout)) * time.Millisecond,
		TLSConfig: &tls.Config{
			ClientAuth: clientAuth,
			ClientCAs:  rootCAs,
			RootCAs:    rootCAs,
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				cert := verifiedChains[0][0]
				log.L(ctx).Debugf("Client certificate provided Subject=%s Issuer=%s Expiry=%s", cert.Subject, cert.Issuer, cert.NotAfter)
				return nil
			},
		},
		ConnContext: func(newCtx context.Context, c net.Conn) context.Context {
			l := log.L(ctx).WithField("req", fftypes.ShortID())
			newCtx = log.WithLogger(newCtx, l)
			l.Debugf("New HTTP connection: remote=%s local=%s", c.RemoteAddr().String(), c.LocalAddr().String())
			return newCtx
		},
	}
	return srv, nil
}

func serveHTTP(ctx context.Context, listener net.Listener, srv *http.Server) (err error) {
	serverEnded := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			log.L(ctx).Infof("API server context cancelled - shutting down")
			srv.Close()
		case <-serverEnded:
			return
		}
	}()

	if config.GetBool(config.HttpTLSEnabled) {
		err = srv.ServeTLS(listener, config.GetString(config.HttpTLSCertFile), config.GetString(config.HttpTLSKeyFile))
	} else {
		err = srv.Serve(listener)
	}
	if err == http.ErrServerClosed {
		err = nil
	}
	close(serverEnded)
	log.L(ctx).Infof("API server complete")

	return err
}

func jsonHandler(e engine.Engine, route *Route) http.HandlerFunc {
	// Check the mandatory parts are ok at startup time
	route.JSONInputValue()
	route.JSONOutputValue()
	return apiWrapper(func(res http.ResponseWriter, req *http.Request) (int, error) {
		l := log.L(req.Context())
		input := route.JSONInputValue()
		var output interface{}
		contentType := req.Header.Get("Content-Type")
		if req.Method != http.MethodGet && !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
			return 415, i18n.NewError(req.Context(), i18n.MsgInvalidContentType)
		}
		var err error
		var status = 400 // if fail parsing input
		if err == nil {
			if input != nil {
				err = json.NewDecoder(req.Body).Decode(&input)
			}
		}
		if err == nil {
			output, status, err = route.JSONHandler(e, req, input)
		}
		if output == nil && err == nil && status != 204 {
			l.Errorf("Request return non-204 code [%d] with no error, and no data", status)
			err = i18n.NewError(req.Context(), i18n.MsgNilResponseNon204)
		}
		if err == nil {
			res.Header().Add("Content-Type", "application/json")
			res.WriteHeader(status)
			if output != nil {
				err = json.NewEncoder(res).Encode(output)
				if err != nil {
					err = i18n.WrapError(req.Context(), err, i18n.MsgResponseMarshalError)
					l.Errorf(err.Error())
				}
			}
		}
		return status, err
	})
}

func apiWrapper(handler func(res http.ResponseWriter, req *http.Request) (status int, err error)) http.HandlerFunc {
	apiTimeout := config.GetUint(config.APIRequestTimeout) // Query once at startup when wrapping
	return func(res http.ResponseWriter, req *http.Request) {

		// Configure a server-side timeout on each request, to try and avoid cases where the API requester
		// times out, and we continue to churn indefinitely processing the request.
		// Long-running processes should be dispatched asynchronously (API returns 202 Accepted asap),
		// and the caller can either listen on the websocket for updates, or poll the status of the affected object.
		// This is dependent on the context being passed down through to all blocking operations down the stack
		// (while avoiding passing the context to asynchronous tasks that are dispatched as a result of the request)
		ctx, cancel := context.WithTimeout(req.Context(), time.Duration(apiTimeout)*time.Millisecond)
		req = req.WithContext(ctx)
		defer cancel()

		// Wrap the request itself in a log wrapper, that gives minimal request/response and timing info
		l := log.L(ctx)
		l.Infof("--> %s %s", req.Method, req.URL.Path)
		startTime := time.Now()
		status, err := handler(res, req)
		duration := float64(time.Since(startTime)) / float64(time.Millisecond)
		if err != nil {
			// Routers don't need to tweak the status code when sending errors.
			// .. either the FF12345 error they raise is mapped to a status hint
			ffcodeExtract := ffcodeExtractor.FindStringSubmatch(err.Error())
			if len(ffcodeExtract) >= 2 {
				if statusHint, ok := i18n.GetStatusHint(ffcodeExtract[1]); ok {
					status = statusHint
				}
			}
			// ... or we default to 500
			if status < 300 {
				status = 500
			}
			l.Infof("<-- %s %s [%d] (%.2fms): %s", req.Method, req.URL.Path, status, duration, err)
			res.Header().Add("Content-Type", "application/json")
			res.WriteHeader(status)
			_ = json.NewEncoder(res).Encode(&RESTError{
				Error: err.Error(),
			})
		} else {
			l.Infof("<-- %s %s [%d] (%.2fms)", req.Method, req.URL.Path, status, duration)
		}
	}
}

func notFoundHandler(res http.ResponseWriter, req *http.Request) (status int, err error) {
	res.Header().Add("Content-Type", "application/json")
	return 404, i18n.NewError(req.Context(), i18n.Msg404NotFound)
}

func createMuxRouter(e engine.Engine) *mux.Router {
	r := mux.NewRouter()
	for _, route := range Routes {
		if route.JSONHandler != nil {
			r.HandleFunc(fmt.Sprintf("/api/v1/%s", route.Path), jsonHandler(e, route)).
				Methods(route.Method)
		}
	}
	r.NotFoundHandler = apiWrapper(notFoundHandler)
	return r
}
