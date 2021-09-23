package pprof

import (
	"context"
	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/logger"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
	"net/http"
)

import _ "net/http/pprof"

const hostKey = "host"
const portKey = "port"

type PProf struct {
	mValues     model.ConfigMap
	mHttpServer *http.Server
}

func (p *PProf) Name() string {
	return "abesh_pprof"
}

func (p *PProf) Version() string {
	return constant.Version
}

func (p *PProf) Category() string {
	return string(constant.CategoryNetwork)
}

func (p *PProf) ContractId() string {
	return "abesh:pprof"
}

func (p *PProf) GetConfigMap() model.ConfigMap {
	return p.mValues
}

func (p *PProf) SetConfigMap(values model.ConfigMap) error {
	p.mValues = values
	return nil
}

func (p *PProf) New() iface.ICapability {
	return &PProf{}
}

func (p *PProf) Setup() error {
	p.mHttpServer = new(http.Server)
	p.mHttpServer.Addr = p.mValues.String(hostKey, "127.0.0.1") +
		":" + p.mValues.String(portKey, "6060")

	p.mHttpServer.Handler = http.DefaultServeMux
	return nil
}

func (p *PProf) Start(_ context.Context) error {
	logger.L(p.ContractId()).Info("pprof server started at " + p.mHttpServer.Addr)
	if err := p.mHttpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (p *PProf) Stop(ctx context.Context) error {
	if p.mHttpServer != nil {
		return p.mHttpServer.Shutdown(ctx)
	}

	return nil
}

func init() {
	registry.GlobalRegistry().AddCapability(&PProf{})
}
