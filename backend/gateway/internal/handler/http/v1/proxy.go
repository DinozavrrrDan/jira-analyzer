package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/DinozvrrDan/jira-analyzer/backend/gateway/config"
	"github.com/magellon17/logger"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Proxy struct {
	cfg *config.Config
	log *logger.Logger
}

func NewProxy(log *logger.Logger, cfg *config.Config) *Proxy {
	return &Proxy{
		cfg: cfg,
		log: log,
	}
}

func (proxy *Proxy) NewProxy(serveMux *http.ServeMux, apiVersion string) {
	serveMux.HandleFunc(apiVersion+proxy.cfg.Server.ConnectorHTTP.ConnectorPrefix+"/", proxy.connectorProxy)
	serveMux.HandleFunc(apiVersion+proxy.cfg.Server.AnalyticsHTTP.AnalyticsPrefix+"/", proxy.analyticsProxy)
	serveMux.HandleFunc(apiVersion+proxy.cfg.Server.ResourceHTTP.ResourcePrefix+"/", proxy.resourceProxy)
}

func (proxy *Proxy) connectorProxy(w http.ResponseWriter, r *http.Request) {
	connectorHttp, err := url.Parse(fmt.Sprintf("http://%s:%s", proxy.cfg.Server.ConnectorHTTP.ConnectorHost,
		proxy.cfg.Server.ConnectorHTTP.ConnectorPort))

	if err != nil {
		errorWriter(w, proxy.log, err.Error(), http.StatusBadRequest)
		return
	}

	proxyConnector := httputil.NewSingleHostReverseProxy(connectorHttp)

	proxyConnector.ServeHTTP(w, r)

}

func (proxy *Proxy) analyticsProxy(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(proxy.cfg.Timeout.Analytics))
	defer cancel()

	connectorHttp, err := url.Parse(fmt.Sprintf("http://%s:%s", proxy.cfg.Server.AnalyticsHTTP.AnalyticsHost,
		proxy.cfg.Server.AnalyticsHTTP.AnalyticsPort))

	if err != nil {
		errorWriter(w, proxy.log, err.Error(), http.StatusBadRequest)
		return
	}

	r = r.WithContext(ctx)

	proxyAnalytics := httputil.NewSingleHostReverseProxy(connectorHttp)

	proxyAnalytics.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if ctxStatus := r.Context().Err(); errors.Is(ctxStatus, context.DeadlineExceeded) {
			http.Error(w, "Request Timeout", http.StatusRequestTimeout)
		} else {
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
		}
	}

	proxyAnalytics.ServeHTTP(w, r)

}

func (proxy *Proxy) resourceProxy(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(proxy.cfg.Timeout.Resource))
	defer cancel()

	resourceHttp, err := url.Parse(fmt.Sprintf("http://%s:%s", proxy.cfg.Server.ResourceHTTP.ResourceHost,
		proxy.cfg.Server.ResourceHTTP.ResourcePort))

	if err != nil {
		errorWriter(w, proxy.log, err.Error(), http.StatusBadRequest)
		return
	}

	r = r.WithContext(ctx)

	proxyResource := httputil.NewSingleHostReverseProxy(resourceHttp)

	proxyResource.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if ctxStatus := r.Context().Err(); errors.Is(ctxStatus, context.DeadlineExceeded) {
			http.Error(w, "Request Timeout", http.StatusRequestTimeout)
		} else {
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
		}
	}

	proxyResource.ServeHTTP(w, r)

}
