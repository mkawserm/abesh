package exrpc

import (
	"context"
	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/logger"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
	"go.uber.org/zap"
	"net/http"
)

type authorizerAndExpression struct {
	authorizer iface.IAuthorizer
	expression string
}

type ExRPC struct {
	mValues           model.ConfigMap
	mHost             string
	mPort             string
	mHttpServer       *http.Server
	mHttpServerMux    *http.ServeMux
	mEventTransmitter iface.IEventTransmitter

	mAuthorizerMap map[string]authorizerAndExpression
}

func (e *ExRPC) Name() string {
	return "abesh_example_rpc"
}

func (e *ExRPC) Version() string {
	return "0.0.1"
}

func (e *ExRPC) Category() string {
	return string(constant.CategoryRPC)
}

func (e *ExRPC) ContractId() string {
	return "abesh:ex_rpc"
}

func (e *ExRPC) GetConfigMap() model.ConfigMap {
	return e.mValues
}

func (e *ExRPC) SetEventTransmitter(eventTransmitter iface.IEventTransmitter) error {
	e.mEventTransmitter = eventTransmitter
	return nil
}

func (e *ExRPC) GetEventTransmitter() iface.IEventTransmitter {
	return e.mEventTransmitter
}

func (e *ExRPC) Setup() error {
	e.mHttpServer = new(http.Server)
	e.mHttpServerMux = new(http.ServeMux)
	e.mAuthorizerMap = make(map[string]authorizerAndExpression)

	// setup server details
	e.mHttpServer.Handler = e.mHttpServerMux
	e.mHttpServer.Addr = e.mHost + ":" + e.mPort

	// setup rpc method
	// /test.TestRPC/Allow

	e.mHttpServerMux.HandleFunc("/test.TestRPC/Allow", func(writer http.ResponseWriter, request *http.Request) {
		ae := e.mAuthorizerMap["/test.TestRPC/Allow"]

		if ae.authorizer == nil {
			writer.WriteHeader(403)
			_, _ = writer.Write([]byte("unauthorized"))
			return
		}

		if ae.authorizer.IsAuthorized(ae.expression, nil) {
			writer.WriteHeader(200)
			_, _ = writer.Write([]byte("allowed"))
			return
		}

		writer.WriteHeader(403)
		_, _ = writer.Write([]byte("unauthorized"))
		return
	})

	e.mHttpServerMux.HandleFunc("/test.TestRPC/Deny", func(writer http.ResponseWriter, request *http.Request) {
		ae := e.mAuthorizerMap["/test.TestRPC/Deny"]

		if ae.authorizer == nil {
			writer.WriteHeader(403)
			_, _ = writer.Write([]byte("unauthorized"))
			return
		}

		if ae.authorizer.IsAuthorized(ae.expression, nil) {
			writer.WriteHeader(200)
			_, _ = writer.Write([]byte("allowed"))
			return
		} else {
			writer.WriteHeader(403)
			_, _ = writer.Write([]byte("denied"))
			return
		}
	})

	return nil
}

func (e *ExRPC) SetConfigMap(values model.ConfigMap) error {
	e.mValues = values

	e.mHost = e.mValues.String("host", "0.0.0.0")
	e.mPort = e.mValues.String("port", "8081")

	return nil
}

func (e *ExRPC) New() iface.ICapability {
	return &ExRPC{}
}

func (e *ExRPC) Start(_ context.Context) error {
	logger.L(e.ContractId()).Info("ex rpc server started at " + e.mHttpServer.Addr)
	if err := e.mHttpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (e *ExRPC) Stop(ctx context.Context) error {
	if e.mHttpServer != nil {
		return e.mHttpServer.Shutdown(ctx)
	}

	return nil
}

func (e *ExRPC) TransmitInputEvent(contractId string, inputEvent *model.Event) {
	if e.GetEventTransmitter() != nil {
		go func() {
			err := e.GetEventTransmitter().TransmitInputEvent(contractId, inputEvent)
			if err != nil {
				logger.L(e.ContractId()).Error(err.Error(),
					zap.String("version", e.Version()),
					zap.String("name", e.Name()),
					zap.String("contract_id", e.ContractId()))
			}

		}()
	}
}

func (e *ExRPC) TransmitOutputEvent(contractId string, outputEvent *model.Event) {
	if e.GetEventTransmitter() != nil {
		go func() {
			err := e.GetEventTransmitter().TransmitOutputEvent(contractId, outputEvent)
			if err != nil {
				logger.L(e.ContractId()).Error(err.Error(),
					zap.String("version", e.Version()),
					zap.String("name", e.Name()),
					zap.String("contract_id", e.ContractId()))
			}
		}()
	}
}

func (e *ExRPC) AddAuthorizer(authorizer iface.IAuthorizer, authorizerExpression string, method string) error {
	e.mAuthorizerMap[method] = authorizerAndExpression{authorizer: authorizer, expression: authorizerExpression}
	return nil
}

func init() {
	registry.GlobalRegistry().AddCapability(&ExRPC{})
}
