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

func GetHeader(headers map[string]string, key string) string {
	value, ok := headers[key]
	if ok {
		return value
	}

	return ""
}

type HTTPServer struct {
	mHost   string
	mPort   string
	mValues map[string]string

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

	logger.L(constant.Name).Info("http server setup complete",
		zap.String("host", h.mHost),
		zap.String("port", h.mPort))

	return nil
}

func (h *HTTPServer) Start(_ context.Context) error {
	logger.L(constant.Name).Info("http server started at " + h.mHttpServer.Addr)
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
		timerStart := time.Now()
		defer func() {
			logger.L(h.ContractId()).Debug("request completed")
			elapsed := time.Since(timerStart)
			logger.L(constant.Name).Debug("request execution time", zap.Duration("seconds", elapsed))
		}()

		logger.L(h.ContractId()).Debug("request stated")
		logger.L(h.ContractId()).Debug("request data",
			zap.String("path", request.URL.Path),
			zap.String("method", request.Method),
			zap.String("path_with_query", request.RequestURI))

		if method != request.Method {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var data []byte
		var err error

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
				writer.WriteHeader(http.StatusForbidden)
				return
			}
		}

		if data, err = ioutil.ReadAll(request.Body); err != nil {
			logger.S(constant.Name).Error(err.Error(),
				zap.String("version", h.Version()),
				zap.String("name", h.Name()),
				zap.String("contract_id", h.ContractId()))

			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		inputEvent := &model.Event{
			Metadata: metadata,
			TypeUrl:  GetHeader(headers, "content-type"),
			Value:    data,
		}

		var outputEvent *model.Event

		outputEvent, err = service.Serve(request.Context(), capabilityRegistry, inputEvent)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		//NOTE: handle success
		writer.WriteHeader(int(outputEvent.Metadata.StatusCode))
		for k, v := range outputEvent.Metadata.Headers {
			writer.Header().Add(k, v)
		}

		if _, err = writer.Write(outputEvent.Value); err != nil {
			logger.S(constant.Name).Error(err.Error(),
				zap.String("version", h.Version()),
				zap.String("name", h.Name()),
				zap.String("contract_id", h.ContractId()))
		}
	})

	return nil
}

func init() {
	registry.GlobalRegistry().AddCapability(&HTTPServer{})
}
