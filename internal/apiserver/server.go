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

package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hyperledger/firefly-common/pkg/ffapi"
	"github.com/hyperledger/firefly/internal/coreconfig"
	"github.com/hyperledger/firefly/internal/coremsgs"
	"github.com/hyperledger/firefly/internal/metrics"
	"github.com/hyperledger/firefly/internal/namespace"
	"github.com/hyperledger/firefly/pkg/database"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ghodss/yaml"
	"github.com/gorilla/mux"

	"github.com/hyperledger/firefly-common/pkg/config"
	"github.com/hyperledger/firefly-common/pkg/httpserver"
	"github.com/hyperledger/firefly-common/pkg/i18n"
	"github.com/hyperledger/firefly/internal/events/eifactory"
	"github.com/hyperledger/firefly/internal/events/websockets"
)

var (
	adminConfig   = config.RootSection("admin")
	apiConfig     = config.RootSection("http")
	metricsConfig = config.RootSection("metrics")
	corsConfig    = config.RootSection("cors")
)

// Server is the external interface for the API Server
type Server interface {
	Serve(ctx context.Context, mgr namespace.Manager) error
}

type apiServer struct {
	// Defaults set with config
	defaultFilterLimit uint64
	maxFilterLimit     uint64
	maxFilterSkip      uint64
	apiTimeout         time.Duration
	apiMaxTimeout      time.Duration
	metricsEnabled     bool
	ffiSwaggerGen      FFISwaggerGen
}

func InitConfig() {
	httpserver.InitHTTPConfig(apiConfig, 5000)
	httpserver.InitHTTPConfig(adminConfig, 5001)
	httpserver.InitHTTPConfig(metricsConfig, 6000)
	httpserver.InitCORSConfig(corsConfig)
	initMetricsConfig(metricsConfig)
}

func NewAPIServer() Server {
	return &apiServer{
		defaultFilterLimit: uint64(config.GetUint(coreconfig.APIDefaultFilterLimit)),
		maxFilterLimit:     uint64(config.GetUint(coreconfig.APIMaxFilterLimit)),
		maxFilterSkip:      uint64(config.GetUint(coreconfig.APIMaxFilterSkip)),
		apiTimeout:         config.GetDuration(coreconfig.APIRequestTimeout),
		apiMaxTimeout:      config.GetDuration(coreconfig.APIRequestMaxTimeout),
		metricsEnabled:     config.GetBool(coreconfig.MetricsEnabled),
		ffiSwaggerGen:      NewFFISwaggerGen(),
	}
}

// Serve is the main entry point for the API Server
func (as *apiServer) Serve(ctx context.Context, mgr namespace.Manager) (err error) {
	httpErrChan := make(chan error)
	adminErrChan := make(chan error)
	metricsErrChan := make(chan error)

	apiHTTPServer, err := httpserver.NewHTTPServer(ctx, "api", as.createMuxRouter(ctx, mgr), httpErrChan, apiConfig, corsConfig)
	if err != nil {
		return err
	}
	go apiHTTPServer.ServeHTTP(ctx)

	if config.GetBool(coreconfig.AdminEnabled) {
		adminHTTPServer, err := httpserver.NewHTTPServer(ctx, "admin", as.createAdminMuxRouter(mgr), adminErrChan, adminConfig, corsConfig)
		if err != nil {
			return err
		}
		go adminHTTPServer.ServeHTTP(ctx)
	}

	if as.metricsEnabled {
		metricsHTTPServer, err := httpserver.NewHTTPServer(ctx, "metrics", as.createMetricsMuxRouter(), metricsErrChan, metricsConfig, corsConfig)
		if err != nil {
			return err
		}
		go metricsHTTPServer.ServeHTTP(ctx)
	}

	return as.waitForServerStop(httpErrChan, adminErrChan, metricsErrChan)
}

func (as *apiServer) waitForServerStop(httpErrChan, adminErrChan, metricsErrChan chan error) error {
	select {
	case err := <-httpErrChan:
		return err
	case err := <-adminErrChan:
		return err
	case err := <-metricsErrChan:
		return err
	}
}

func (as *apiServer) getPublicURL(conf config.Section, pathPrefix string) string {
	publicURL := conf.GetString(httpserver.HTTPConfPublicURL)
	if publicURL == "" {
		proto := "https"
		if !conf.GetBool(httpserver.HTTPConfTLSEnabled) {
			proto = "http"
		}
		publicURL = fmt.Sprintf("%s://%s:%s", proto, conf.GetString(httpserver.HTTPConfAddress), conf.GetString(httpserver.HTTPConfPort))
	}
	if pathPrefix != "" {
		publicURL += "/" + pathPrefix
	}
	return publicURL
}

func (as *apiServer) swaggerGenConf(apiBaseURL string) *ffapi.Options {
	return &ffapi.Options{
		BaseURL:                   apiBaseURL,
		Title:                     "FireFly",
		Version:                   "1.0",
		PanicOnMissingDescription: config.GetBool(coreconfig.APIOASPanicOnMissingDescription),
		DefaultRequestTimeout:     config.GetDuration(coreconfig.APIRequestTimeout),
		RouteCustomizations: func(ctx context.Context, sg *ffapi.SwaggerGen, route *ffapi.Route, op *openapi3.Operation) {
			if ce, ok := route.Extensions.(*coreExtensions); ok {
				if ce.FilterFactory != nil {
					fields := ce.FilterFactory.NewFilter(ctx).Fields()
					sort.Strings(fields)
					for _, field := range fields {
						sg.AddParam(ctx, op, "query", field, "", "", coremsgs.APIFilterParamDesc, false)
					}
					sg.AddParam(ctx, op, "query", "sort", "", "", coremsgs.APIFilterSortDesc, false)
					sg.AddParam(ctx, op, "query", "ascending", "", "", coremsgs.APIFilterAscendingDesc, false)
					sg.AddParam(ctx, op, "query", "descending", "", "", coremsgs.APIFilterDescendingDesc, false)
					sg.AddParam(ctx, op, "query", "skip", "", "", coremsgs.APIFilterSkipDesc, false, config.GetUint(coreconfig.APIMaxFilterSkip))
					sg.AddParam(ctx, op, "query", "limit", "", config.GetString(coreconfig.APIDefaultFilterLimit), coremsgs.APIFilterLimitDesc, false, config.GetUint(coreconfig.APIMaxFilterLimit))
					sg.AddParam(ctx, op, "query", "count", "", "", coremsgs.APIFilterCountDesc, false)
				}
			}
		},
	}
}

func (as *apiServer) swaggerHandler(generator func(req *http.Request) (*openapi3.T, error)) func(res http.ResponseWriter, req *http.Request) (status int, err error) {
	return func(res http.ResponseWriter, req *http.Request) (status int, err error) {
		vars := mux.Vars(req)
		doc, err := generator(req)
		if err != nil {
			return 500, err
		}
		if vars["ext"] == ".json" {
			res.Header().Add("Content-Type", "application/json")
			b, _ := json.Marshal(&doc)
			_, _ = res.Write(b)
		} else {
			res.Header().Add("Content-Type", "application/x-yaml")
			b, _ := yaml.Marshal(&doc)
			_, _ = res.Write(b)
		}
		return 200, nil
	}
}

func (as *apiServer) swaggerGenerator(routes []*ffapi.Route, apiBaseURL string) func(req *http.Request) (*openapi3.T, error) {
	swg := ffapi.NewSwaggerGen(as.swaggerGenConf(apiBaseURL))
	return func(req *http.Request) (*openapi3.T, error) {
		return swg.Generate(req.Context(), routes), nil
	}
}

func (as *apiServer) contractSwaggerGenerator(mgr namespace.Manager, apiBaseURL string) func(req *http.Request) (*openapi3.T, error) {
	return func(req *http.Request) (*openapi3.T, error) {
		vars := mux.Vars(req)
		cm := mgr.Orchestrator(vars["ns"]).Contracts()
		api, err := cm.GetContractAPI(req.Context(), apiBaseURL, vars["ns"], vars["apiName"])
		if err != nil {
			return nil, err
		} else if api == nil || api.Interface == nil {
			return nil, i18n.NewError(req.Context(), coremsgs.Msg404NoResult)
		}

		ffi, err := cm.GetFFIByIDWithChildren(req.Context(), api.Interface.ID)
		if err != nil {
			return nil, err
		}

		baseURL := fmt.Sprintf("%s/namespaces/%s/apis/%s", apiBaseURL, vars["ns"], vars["apiName"])
		return as.ffiSwaggerGen.Generate(req.Context(), baseURL, api, ffi), nil
	}
}

func (as *apiServer) routeHandler(hf *ffapi.HandlerFactory, mgr namespace.Manager, apiBaseURL string, route *ffapi.Route) http.HandlerFunc {
	// We extend the base ffapi functionality, with standardized DB filter support for all core resources.
	// We also pass the Orchestrator context through
	ce := route.Extensions.(*coreExtensions)
	route.JSONHandler = func(r *ffapi.APIRequest) (output interface{}, err error) {
		var filter database.AndFilter
		if ce.FilterFactory != nil {
			filter, err = as.buildFilter(r.Req, ce.FilterFactory)
			if err != nil {
				return nil, err
			}
		}
		vars := mux.Vars(r.Req)
		or := mgr.Orchestrator(extractNamespace(vars))
		cr := &coreRequest{
			mgr:        mgr,
			or:         or,
			ctx:        r.Req.Context(),
			filter:     filter,
			apiBaseURL: apiBaseURL,
		}
		return ce.CoreJSONHandler(r, cr)
	}
	if ce.CoreFormUploadHandler != nil {
		route.FormUploadHandler = func(r *ffapi.APIRequest) (output interface{}, err error) {
			vars := mux.Vars(r.Req)
			or := mgr.Orchestrator(vars["ns"])
			cr := &coreRequest{
				or:         or,
				ctx:        r.Req.Context(),
				apiBaseURL: apiBaseURL,
			}
			return ce.CoreFormUploadHandler(r, cr)
		}
	}
	return hf.RouteHandler(route)
}

func (as *apiServer) handlerFactory() *ffapi.HandlerFactory {
	return &ffapi.HandlerFactory{
		DefaultRequestTimeout: config.GetDuration(coreconfig.APIRequestTimeout),
		MaxTimeout:            config.GetDuration(coreconfig.APIRequestMaxTimeout),
	}
}

func (as *apiServer) createMuxRouter(ctx context.Context, mgr namespace.Manager) *mux.Router {
	r := mux.NewRouter()
	hf := as.handlerFactory()

	if as.metricsEnabled {
		r.Use(metrics.GetRestServerInstrumentation().Middleware)
	}

	publicURL := as.getPublicURL(apiConfig, "")
	apiBaseURL := fmt.Sprintf("%s/api/v1", publicURL)
	for _, route := range routes {
		if ce, ok := route.Extensions.(*coreExtensions); ok {
			if ce.CoreJSONHandler != nil {
				r.HandleFunc(fmt.Sprintf("/api/v1/%s", route.Path), as.routeHandler(hf, mgr, apiBaseURL, route)).
					Methods(route.Method)
			}
		}
	}

	r.HandleFunc(`/api/v1/namespaces/{ns}/apis/{apiName}/api/swagger{ext:\.yaml|\.json|}`, hf.APIWrapper(as.swaggerHandler(as.contractSwaggerGenerator(mgr, apiBaseURL))))
	r.HandleFunc(`/api/v1/namespaces/{ns}/apis/{apiName}/api`, func(rw http.ResponseWriter, req *http.Request) {
		url := req.URL.String() + "/swagger.yaml"
		handler := hf.APIWrapper(hf.SwaggerUIHandler(url))
		handler(rw, req)
	})

	r.HandleFunc(`/api/swagger{ext:\.yaml|\.json|}`, hf.APIWrapper(as.swaggerHandler(as.swaggerGenerator(routes, apiBaseURL))))
	r.HandleFunc(`/api`, hf.APIWrapper(hf.SwaggerUIHandler(publicURL+"/api/swagger.yaml")))
	r.HandleFunc(`/favicon{any:.*}.png`, favIcons)

	ws, _ := eifactory.GetPlugin(ctx, "websockets")
	r.HandleFunc(`/ws`, ws.(*websockets.WebSockets).ServeHTTP)

	uiPath := config.GetString(coreconfig.UIPath)
	if uiPath != "" && config.GetBool(coreconfig.UIEnabled) {
		r.PathPrefix(`/ui`).Handler(newStaticHandler(uiPath, "index.html", `/ui`))
	}

	r.NotFoundHandler = hf.APIWrapper(as.notFoundHandler)
	return r
}

func (as *apiServer) notFoundHandler(res http.ResponseWriter, req *http.Request) (status int, err error) {
	res.Header().Add("Content-Type", "application/json")
	return 404, i18n.NewError(req.Context(), coremsgs.Msg404NotFound)
}

func (as *apiServer) adminWSHandler(mgr namespace.Manager) http.HandlerFunc {
	// The admin events listener will be initialized when we start, so we access it it from Orchestrator on demand
	return func(w http.ResponseWriter, r *http.Request) {
		mgr.AdminEvents().ServeHTTPWebSocketListener(w, r)
	}
}

func (as *apiServer) createAdminMuxRouter(mgr namespace.Manager) *mux.Router {
	r := mux.NewRouter()
	if as.metricsEnabled {
		r.Use(metrics.GetAdminServerInstrumentation().Middleware)
	}
	hf := as.handlerFactory()

	publicURL := as.getPublicURL(adminConfig, "admin")
	apiBaseURL := fmt.Sprintf("%s/admin/api/v1", publicURL)
	for _, route := range adminRoutes {
		if ce, ok := route.Extensions.(*coreExtensions); ok {
			if ce.CoreJSONHandler != nil {
				r.HandleFunc(fmt.Sprintf("/admin/api/v1/%s", route.Path), as.routeHandler(hf, mgr, apiBaseURL, route)).
					Methods(route.Method)
			}
		}
	}
	r.HandleFunc(`/admin/api/swagger{ext:\.yaml|\.json|}`, hf.APIWrapper(as.swaggerHandler(as.swaggerGenerator(adminRoutes, apiBaseURL))))
	r.HandleFunc(`/admin/api`, hf.APIWrapper(hf.SwaggerUIHandler(publicURL+"/api/swagger.yaml")))
	r.HandleFunc(`/favicon{any:.*}.png`, favIcons)

	r.HandleFunc(`/admin/ws`, as.adminWSHandler(mgr))

	return r
}

func (as *apiServer) createMetricsMuxRouter() *mux.Router {
	r := mux.NewRouter()

	r.Path(config.GetString(coreconfig.MetricsPath)).Handler(promhttp.InstrumentMetricHandler(metrics.Registry(),
		promhttp.HandlerFor(metrics.Registry(), promhttp.HandlerOpts{})))

	return r
}
