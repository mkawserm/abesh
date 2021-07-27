package httpserver

import (
	"context"
	"errors"
	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/logger"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var ErrPathNotDefined = errors.New("path not defined")
var ErrMethodNotDefined = errors.New("method not defined")

type EventResponse struct {
	Error error
	Event *model.Event
}

func GetValue(headers map[string]string, key string, defaultValue string) string {
	value, ok := headers[key]
	if ok {
		return value
	}

	return defaultValue
}

type HTTPServer struct {
	mHost           string
	mPort           string
	mRequestTimeout time.Duration
	mValues         map[string]string

	mHttpServer    *http.Server
	mHttpServerMux *http.ServeMux
}

func (h *HTTPServer) Name() string {
	return "golang_net_http_server"
}

func (h *HTTPServer) Version() string {
	return "0.0.1"
}

func (h *HTTPServer) Category() string {
	return string(constant.CategoryTrigger)
}

func (h *HTTPServer) ContractId() string {
	return "abesh:httpserver"
}

func (h *HTTPServer) Values() map[string]string {
	return h.mValues
}

func (h *HTTPServer) SetValues(values map[string]string) error {
	h.mValues = values

	if host, ok := values["host"]; ok {
		h.mHost = host
	} else {
		h.mHost = "0.0.0.0"
	}

	if port, ok := values["port"]; ok {
		h.mPort = port
	} else {
		h.mPort = "8080"
	}

	requestTimeout, err := time.ParseDuration(GetValue(h.mValues, "default_request_timeout", "60s"))
	if err != nil {
		return err
	}

	h.mRequestTimeout = requestTimeout
	return nil
}

func (h *HTTPServer) New() iface.ICapability {
	return &HTTPServer{}
}

func (h *HTTPServer) Setup() error {
	h.mHttpServer = new(http.Server)
	h.mHttpServerMux = new(http.ServeMux)

	// setup server details
	h.mHttpServer.Handler = h.mHttpServerMux
	h.mHttpServer.Addr = h.mHost + ":" + h.mPort

	logger.L(h.ContractId()).Info("http server setup complete",
		zap.String("host", h.mHost),
		zap.String("port", h.mPort))

	return nil
}

func (h *HTTPServer) Start(_ context.Context) error {
	logger.L(h.ContractId()).Info("http server started at " + h.mHttpServer.Addr)
	if err := h.mHttpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (h *HTTPServer) Stop(ctx context.Context) error {
	if h.mHttpServer != nil {
		return h.mHttpServer.Shutdown(ctx)
	}

	return nil
}

func (h *HTTPServer) AddService(
	authorizationHandler iface.AuthorizationHandler,
	authorizationExpression string,
	triggerValues map[string]string,
	capabilityRegistry iface.ICapabilityRegistry,
	service iface.IService) error {

	var method string
	var path string
	var ok bool

	if method, ok = triggerValues["method"]; !ok {
		return ErrMethodNotDefined
	}

	method = strings.ToUpper(strings.TrimSpace(method))

	if path, ok = triggerValues["path"]; !ok {
		return ErrPathNotDefined
	}

	path = strings.TrimSpace(path)

	h.mHttpServerMux.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		var err error
		timerStart := time.Now()

		defer func() {
			logger.L(h.ContractId()).Debug("request completed")
			elapsed := time.Since(timerStart)
			logger.L(h.ContractId()).Debug("request execution time", zap.Duration("seconds", elapsed))
		}()

		logger.L(h.ContractId()).Debug("request timeout", zap.Duration("timeout", h.mRequestTimeout))

		logger.L(h.ContractId()).Debug("request stated")
		logger.L(h.ContractId()).Debug("request data",
			zap.String("path", request.URL.Path),
			zap.String("method", request.Method),
			zap.String("path_with_query", request.RequestURI))

		if method != request.Method {
			writer.Header().Add("Content-Type", GetValue(h.mValues, "default_content_type", "application/text"))
			writer.WriteHeader(http.StatusMethodNotAllowed)
			if _, err = writer.Write([]byte(GetValue(h.mValues, "s405m", "405 ERROR"))); err != nil {
				logger.S(h.ContractId()).Error(err.Error(),
					zap.String("version", h.Version()),
					zap.String("name", h.Name()),
					zap.String("contract_id", h.ContractId()))
			}
			return
		}

		var data []byte

		headers := make(map[string]string)

		metadata := &model.Metadata{}
		metadata.Method = request.Method
		metadata.Path = request.URL.EscapedPath()
		metadata.Headers = make(map[string]string)
		metadata.Query = make(map[string]string)

		for k, v := range request.Header {
			if len(v) > 0 {
				metadata.Headers[k] = v[0]
				headers[strings.ToLower(strings.TrimSpace(k))] = v[0]
			}
		}

		for k, v := range request.URL.Query() {
			if len(v) > 0 {
				metadata.Query[k] = v[0]
			}
		}

		if authorizationHandler != nil {
			if !authorizationHandler(authorizationExpression, metadata) {
				writer.Header().Add("Content-Type", GetValue(h.mValues, "default_content_type", "application/text"))
				writer.WriteHeader(http.StatusForbidden)
				if _, err = writer.Write([]byte(GetValue(h.mValues, "s403m", "403 ERROR"))); err != nil {
					logger.S(h.ContractId()).Error(err.Error(),
						zap.String("version", h.Version()),
						zap.String("name", h.Name()),
						zap.String("contract_id", h.ContractId()))
				}
				return
			}
		}

		if data, err = ioutil.ReadAll(request.Body); err != nil {
			logger.S(h.ContractId()).Error(err.Error(),
				zap.String("version", h.Version()),
				zap.String("name", h.Name()),
				zap.String("contract_id", h.ContractId()))

			writer.Header().Add("Content-Type", GetValue(h.mValues, "default_content_type", "application/text"))
			writer.WriteHeader(http.StatusInternalServerError)
			if _, err = writer.Write([]byte(GetValue(h.mValues, "s500m", "500 ERROR"))); err != nil {
				logger.S(h.ContractId()).Error(err.Error(),
					zap.String("version", h.Version()),
					zap.String("name", h.Name()),
					zap.String("contract_id", h.ContractId()))
			}
			return
		}

		inputEvent := &model.Event{
			Metadata: metadata,
			TypeUrl:  GetValue(headers, "content-type", "application/text"),
			Value:    data,
		}

		nCtx, cancel := context.WithTimeout(request.Context(), h.mRequestTimeout)
		defer cancel()

		ch := make(chan EventResponse, 1)

		func() {
			if request.Context().Err() != nil {
				ch <- EventResponse{
					Event: nil,
					Error: request.Context().Err(),
				}
			} else {
				go func() {
					event, errInner := service.Serve(nCtx, capabilityRegistry, inputEvent)
					ch <- EventResponse{Event: event, Error: errInner}
				}()
			}
		}()

		select {
		case <-nCtx.Done():
			writer.Header().Add("Content-Type", GetValue(h.mValues, "default_content_type", "application/text"))
			writer.WriteHeader(http.StatusRequestTimeout)
			if _, err = writer.Write([]byte(GetValue(h.mValues, "s408m", "408 ERROR"))); err != nil {
				logger.S(h.ContractId()).Error(err.Error(),
					zap.String("version", h.Version()),
					zap.String("name", h.Name()),
					zap.String("contract_id", h.ContractId()))
			}

			return
		case r := <-ch:
			if r.Error == context.DeadlineExceeded {
				logger.S(h.ContractId()).Error(r.Error.Error(),
					zap.String("version", h.Version()),
					zap.String("name", h.Name()),
					zap.String("contract_id", h.ContractId()))

				writer.Header().Add("Content-Type", GetValue(h.mValues, "default_content_type", "application/text"))
				writer.WriteHeader(http.StatusRequestTimeout)

				if _, err = writer.Write([]byte(GetValue(h.mValues, "s408m", "408 ERROR"))); err != nil {
					logger.S(h.ContractId()).Error(err,
						zap.String("version", h.Version()),
						zap.String("name", h.Name()),
						zap.String("contract_id", h.ContractId()))
				}
				return
			}

			if r.Error == context.Canceled {
				logger.S(h.ContractId()).Error(r.Error.Error(),
					zap.String("version", h.Version()),
					zap.String("name", h.Name()),
					zap.String("contract_id", h.ContractId()))

				writer.Header().Add("Content-Type", GetValue(h.mValues, "default_content_type", "application/text"))
				writer.WriteHeader(499)

				if _, err = writer.Write([]byte(GetValue(h.mValues, "s499m", "499 ERROR"))); err != nil {
					logger.S(h.ContractId()).Error(err,
						zap.String("version", h.Version()),
						zap.String("name", h.Name()),
						zap.String("contract_id", h.ContractId()))
				}
				return
			}

			if r.Error != nil {
				logger.S(h.ContractId()).Error(r.Error.Error(),
					zap.String("version", h.Version()),
					zap.String("name", h.Name()),
					zap.String("contract_id", h.ContractId()))

				writer.Header().Add("Content-Type", GetValue(h.mValues, "default_content_type", "application/text"))
				writer.WriteHeader(http.StatusInternalServerError)
				if _, err = writer.Write([]byte(GetValue(h.mValues, "s500m", "500 ERROR"))); err != nil {
					logger.S(h.ContractId()).Error(err.Error(),
						zap.String("version", h.Version()),
						zap.String("name", h.Name()),
						zap.String("contract_id", h.ContractId()))
				}
				return
			}

			//NOTE: handle success from service
			for k, v := range r.Event.Metadata.Headers {
				writer.Header().Add(k, v)
			}
			writer.WriteHeader(int(r.Event.Metadata.StatusCode))

			if _, err = writer.Write(r.Event.Value); err != nil {
				logger.S(h.ContractId()).Error(err.Error(),
					zap.String("version", h.Version()),
					zap.String("name", h.Name()),
					zap.String("contract_id", h.ContractId()))
			}
		}
	})

	return nil
}

func init() {
	registry.GlobalRegistry().AddCapability(&HTTPServer{})
}
