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

type HTTPServer struct {
	mHost          string
	mPort          string
	mValues        map[string]string
	mHttpServerMux *http.ServeMux
	mHttpServer    *http.Server
}

func (h *HTTPServer) Name() string {
	return "golang_net_http_server"
}

func (h *HTTPServer) Source() string {
	return "github.com/mkawserm/abesh/capability/httpserver"
}

func (h *HTTPServer) Runtime() string {
	return string(constant.RuntimeNative)
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
	h.mHttpServerMux = new(http.ServeMux)
	h.mHttpServer = new(http.Server)
	return nil
}

func (h *HTTPServer) Start() error {
	h.mHttpServer.Addr = h.mHost + ":" + h.mPort
	h.mHttpServer.Handler = h.mHttpServerMux

	return nil
}

func (h *HTTPServer) Stop() error {
	if h.mHttpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return h.mHttpServer.Shutdown(ctx)
	}

	return nil
}

func (h *HTTPServer) AddService(triggerValues map[string]string, service iface.IService) error {
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
		if method != request.Method {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var data []byte
		var err error

		metadata := &model.Metadata{}
		metadata.Method = request.Method
		metadata.Path = request.RequestURI

		for k, v := range request.Header {
			if len(v) > 0 {
				metadata.Headers[k] = v[0]
			}
		}

		if data, err = ioutil.ReadAll(request.Body); err != nil {
			logger.S(constant.Name).Error(err.Error(), zap.String("contract_id", h.ContractId()))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		inputEvent := &model.Event{
			Metadata: metadata,
			Data:     data,
		}

		var outputEvent *model.Event

		outputEvent, err = service.Process(request.Context(), nil, inputEvent)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		//NOTE: handle success
		writer.WriteHeader(int(outputEvent.Metadata.StatusCode))
		for k, v := range outputEvent.Metadata.Headers {
			writer.Header().Add(k, v)
		}

		if _, err = writer.Write(outputEvent.Data); err != nil {
			logger.S(constant.Name).Error(err.Error(), zap.String("contract_id", h.ContractId()))
		}
	})

	return nil
}

func init() {
	registry.GlobalRegistry().AddCapability(&HTTPServer{})
}
