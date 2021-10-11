package httpserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/logger"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
	"github.com/mkawserm/abesh/utility"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

var ErrPathNotDefined = errors.New("path not defined")
var ErrMethodNotDefined = errors.New("method not defined")

type EventResponse struct {
	Error error
	Event *model.Event
}

type HTTPServer struct {
	mHost                     string
	mPort                     string
	mRequestTimeout           time.Duration
	mDefault404HandlerEnabled bool
	mValues                   model.ConfigMap

	mHttpServer       *http.Server
	mHttpServerMux    *http.ServeMux
	mEventTransmitter iface.IEventTransmitter
}

func (h *HTTPServer) Name() string {
	return "abesh_httpserver"
}

func (h *HTTPServer) Version() string {
	return constant.Version
}

func (h *HTTPServer) Category() string {
	return string(constant.CategoryTrigger)
}

func (h *HTTPServer) ContractId() string {
	return "abesh:httpserver"
}

func (h *HTTPServer) GetConfigMap() model.ConfigMap {
	return h.mValues
}

func (h *HTTPServer) SetConfigMap(values model.ConfigMap) error {
	h.mValues = values
	h.mHost = h.mValues.String("host", "0.0.0.0")
	h.mPort = h.mValues.String("port", "8080")
	h.mRequestTimeout = h.mValues.Duration("default_request_timeout", time.Second)
	h.mDefault404HandlerEnabled = h.mValues.Bool("default_404_handler_enabled", true)
	return nil
}

func (h *HTTPServer) SetEventTransmitter(eventTransmitter iface.IEventTransmitter) error {
	h.mEventTransmitter = eventTransmitter
	return nil
}

func (h *HTTPServer) GetEventTransmitter() iface.IEventTransmitter {
	return h.mEventTransmitter
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

	if h.mDefault404HandlerEnabled {
		h.mHttpServerMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			h.debugMessage(request)

			timerStart := time.Now()

			defer func() {
				logger.L(h.ContractId()).Debug("request completed")
				elapsed := time.Since(timerStart)
				logger.L(h.ContractId()).Debug("request execution time", zap.Duration("seconds", elapsed))
			}()

			h.s404m(writer, nil)
			return
		})
	}

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

func (h *HTTPServer) TransmitInputEvent(contractId string, inputEvent *model.Event) {
	if h.GetEventTransmitter() != nil {
		go func() {
			err := h.GetEventTransmitter().TransmitInputEvent(contractId, inputEvent)
			if err != nil {
				logger.L(h.ContractId()).Error(err.Error(),
					zap.String("version", h.Version()),
					zap.String("name", h.Name()),
					zap.String("contract_id", h.ContractId()))
			}

		}()
	}
}

func (h *HTTPServer) TransmitOutputEvent(contractId string, outputEvent *model.Event) {
	if h.GetEventTransmitter() != nil {
		go func() {
			err := h.GetEventTransmitter().TransmitOutputEvent(contractId, outputEvent)
			if err != nil {
				logger.L(h.ContractId()).Error(err.Error(),
					zap.String("version", h.Version()),
					zap.String("name", h.Name()),
					zap.String("contract_id", h.ContractId()))
			}
		}()
	}
}

func (h *HTTPServer) s403m(writer http.ResponseWriter, errLocal error) {
	if errLocal != nil {
		logger.L(h.ContractId()).Error(errLocal.Error(),
			zap.String("version", h.Version()),
			zap.String("name", h.Name()),
			zap.String("contract_id", h.ContractId()))
	}

	writer.Header().Add("Content-Type", h.mValues.String("default_content_type", "application/text"))
	writer.WriteHeader(http.StatusForbidden)
	if _, err := writer.Write([]byte(h.mValues.String("s403m", "403 ERROR"))); err != nil {
		logger.L(h.ContractId()).Error(err.Error(),
			zap.String("version", h.Version()),
			zap.String("name", h.Name()),
			zap.String("contract_id", h.ContractId()))
	}
}

func (h *HTTPServer) s404m(writer http.ResponseWriter, errLocal error) {
	if errLocal != nil {
		logger.L(h.ContractId()).Error(errLocal.Error(),
			zap.String("version", h.Version()),
			zap.String("name", h.Name()),
			zap.String("contract_id", h.ContractId()))
	}

	writer.Header().Add("Content-Type", h.mValues.String("default_content_type", "application/text"))
	writer.WriteHeader(http.StatusNotFound)
	if _, err := writer.Write([]byte(h.mValues.String("s404m", "404 ERROR"))); err != nil {
		logger.L(h.ContractId()).Error(err.Error(),
			zap.String("version", h.Version()),
			zap.String("name", h.Name()),
			zap.String("contract_id", h.ContractId()))
	}
}

func (h *HTTPServer) s405m(writer http.ResponseWriter, errLocal error) {
	if errLocal != nil {
		logger.L(h.ContractId()).Error(errLocal.Error(),
			zap.String("version", h.Version()),
			zap.String("name", h.Name()),
			zap.String("contract_id", h.ContractId()))
	}

	writer.Header().Add("Content-Type", h.mValues.String("default_content_type", "application/text"))
	writer.WriteHeader(http.StatusMethodNotAllowed)
	if _, err := writer.Write([]byte(h.mValues.String("s405m", "405 ERROR"))); err != nil {
		logger.L(h.ContractId()).Error(err.Error(),
			zap.String("version", h.Version()),
			zap.String("name", h.Name()),
			zap.String("contract_id", h.ContractId()))
	}
}

func (h *HTTPServer) s408m(writer http.ResponseWriter, errLocal error) {
	if errLocal != nil {
		logger.L(h.ContractId()).Error(errLocal.Error(),
			zap.String("version", h.Version()),
			zap.String("name", h.Name()),
			zap.String("contract_id", h.ContractId()))
	}

	writer.Header().Add("Content-Type", h.mValues.String("default_content_type", "application/text"))
	writer.WriteHeader(http.StatusRequestTimeout)
	if _, err := writer.Write([]byte(h.mValues.String("s408m", "408 ERROR"))); err != nil {
		logger.L(h.ContractId()).Error(err.Error(),
			zap.String("version", h.Version()),
			zap.String("name", h.Name()),
			zap.String("contract_id", h.ContractId()))
	}
}

func (h *HTTPServer) s499m(writer http.ResponseWriter, errLocal error) {
	if errLocal != nil {
		logger.L(h.ContractId()).Error(errLocal.Error(),
			zap.String("version", h.Version()),
			zap.String("name", h.Name()),
			zap.String("contract_id", h.ContractId()))
	}

	writer.Header().Add("Content-Type", h.mValues.String("default_content_type", "application/text"))
	writer.WriteHeader(499)

	if _, err := writer.Write([]byte(h.mValues.String("s499m", "499 ERROR"))); err != nil {
		logger.L(h.ContractId()).Error(err.Error(),
			zap.String("version", h.Version()),
			zap.String("name", h.Name()),
			zap.String("contract_id", h.ContractId()))
	}
}

func (h *HTTPServer) s500m(writer http.ResponseWriter, errLocal error) {
	if errLocal != nil {
		logger.L(h.ContractId()).Error(errLocal.Error(),
			zap.String("version", h.Version()),
			zap.String("name", h.Name()),
			zap.String("contract_id", h.ContractId()))
	}

	writer.Header().Add("Content-Type", h.mValues.String("default_content_type", "application/text"))
	writer.WriteHeader(http.StatusInternalServerError)
	if _, err := writer.Write([]byte(h.mValues.String("s500m", "500 ERROR"))); err != nil {
		logger.L(h.ContractId()).Error(err.Error(),
			zap.String("version", h.Version()),
			zap.String("name", h.Name()),
			zap.String("contract_id", h.ContractId()))
	}
}

func (h *HTTPServer) debugMessage(request *http.Request) {
	logger.L(h.ContractId()).Debug("request local timeout in seconds", zap.Duration("timeout", h.mRequestTimeout))
	logger.L(h.ContractId()).Debug("request started")
	logger.L(h.ContractId()).Debug("request data",
		zap.String("path", request.URL.Path),
		zap.String("method", request.Method),
		zap.String("path_with_query", request.RequestURI))
}

func (h *HTTPServer) AddService(
	authorizer iface.IAuthorizer,
	authorizerExpression string,
	triggerValues model.ConfigMap,
	service iface.IService) error {

	logger.L(h.ContractId()).Debug("service add",
		zap.Any("authorizer", authorizer),
		zap.Any("expression", authorizerExpression),
		zap.Any("triggerValues", triggerValues))

	var method string
	var path string
	var methodList []string

	if method = triggerValues.String("method", ""); len(method) == 0 {
		return ErrMethodNotDefined
	}

	method = strings.ToUpper(strings.TrimSpace(method))
	methodList = strings.Split(method, ",")
	if len(methodList) > 0 {
		sort.Strings(methodList)
	}

	if path = triggerValues.String("path", ""); len(path) == 0 {
		return ErrPathNotDefined
	}

	path = strings.TrimSpace(path)

	requestHandler := func(writer http.ResponseWriter, request *http.Request) {
		var err error
		timerStart := time.Now()

		defer func() {
			logger.L(h.ContractId()).Debug("request completed")
			elapsed := time.Since(timerStart)
			logger.L(h.ContractId()).Debug("request execution time", zap.Duration("seconds", elapsed))
		}()

		defer func() {
			if r := recover(); r != nil {
				panicMsg := fmt.Sprintf("%v", r)
				logger.L(h.ContractId()).Info("recovering from panic")

				// add as much information as possible
				logger.L(h.ContractId()).Error("panic data",
					zap.String("host_name", request.URL.Hostname()),
					zap.String("host", request.URL.Host),
					zap.String("path", request.URL.Path),
					zap.String("method", request.Method),
					zap.String("uri", request.RequestURI),
					zap.String("panic_msg", panicMsg))

				h.s500m(writer, fmt.Errorf("%v", r))
				return
			}
		}()

		h.debugMessage(request)

		if !utility.IsIn(methodList, method) {
			h.s405m(writer, nil)
			return
		}

		var data []byte

		headers := make(map[string]string)

		metadata := &model.Metadata{}
		metadata.Method = request.Method
		metadata.Path = request.URL.EscapedPath()
		metadata.Headers = make(map[string]string)
		metadata.Query = make(map[string]string)
		metadata.ContractIdList = append(metadata.ContractIdList, h.ContractId())

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

		if authorizer != nil {
			if !authorizer.IsAuthorized(authorizerExpression, metadata) {
				h.s403m(writer, nil)
				return
			}
		}

		if data, err = ioutil.ReadAll(request.Body); err != nil {
			h.s500m(writer, err)
			return
		}

		inputEvent := &model.Event{
			Metadata: metadata,
			TypeUrl:  utility.GetValue(headers, "content-type", "application/text"),
			Value:    data,
		}

		// transmit input event
		h.TransmitInputEvent(service.ContractId(), inputEvent)

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
					event, errInner := service.Serve(nCtx, inputEvent)
					ch <- EventResponse{Event: event, Error: errInner}
				}()
			}
		}()

		select {
		case <-nCtx.Done():
			h.s408m(writer, nil)
			return
		case r := <-ch:
			if r.Error == context.DeadlineExceeded {
				h.s408m(writer, r.Error)
				return
			}

			if r.Error == context.Canceled {
				h.s499m(writer, r.Error)
				return
			}

			if r.Error != nil {
				h.s500m(writer, r.Error)
				return
			}

			// transmit output event
			h.TransmitOutputEvent(service.ContractId(), r.Event)

			//NOTE: handle success from service
			for k, v := range r.Event.Metadata.Headers {
				writer.Header().Add(k, v)
			}
			writer.WriteHeader(int(r.Event.Metadata.StatusCode))

			if _, err = writer.Write(r.Event.Value); err != nil {
				logger.L(h.ContractId()).Error(err.Error(),
					zap.String("version", h.Version()),
					zap.String("name", h.Name()),
					zap.String("contract_id", h.ContractId()))
			}
		}
	}

	h.mHttpServerMux.HandleFunc(path, requestHandler)

	return nil
}

func init() {
	registry.GlobalRegistry().AddCapability(&HTTPServer{})
}
