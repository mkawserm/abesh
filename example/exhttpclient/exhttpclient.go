package echo

import (
	"context"
	"errors"
	httpclient2 "github.com/mkawserm/abesh/capability/httpclient"
	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
	"github.com/mkawserm/abesh/utility"
	"io/ioutil"
)

var ErrHTTPClientNotFound = errors.New("abesh:httpclient not found")

type ExHttpClient struct {
	mValues     iface.ConfigMap
	mUrl        string
	mHttpClient iface.IHttpClient
}

func (e *ExHttpClient) Name() string {
	return "abesh_example_httpclient"
}

func (e *ExHttpClient) Version() string {
	return "0.0.1"
}

func (e *ExHttpClient) Category() string {
	return string(constant.CategoryService)
}

func (e *ExHttpClient) ContractId() string {
	return "abesh:ex_httpclient"
}

func (e *ExHttpClient) GetConfigMap() iface.ConfigMap {
	return e.mValues
}

func (e *ExHttpClient) Setup() error {
	return nil
}

func (e *ExHttpClient) SetCapabilityRegistry(registry iface.ICapabilityRegistry) error {
	e.mHttpClient = httpclient2.GetHttpClient(registry)
	if e.mHttpClient == nil {
		return ErrHTTPClientNotFound
	}

	return nil
}

func (e *ExHttpClient) SetConfigMap(values iface.ConfigMap) error {
	e.mValues = values

	e.mUrl = utility.GetValue(e.mValues, "url", "https://jsonip.com")
	return nil
}

func (e *ExHttpClient) New() iface.ICapability {
	return &ExHttpClient{}
}

func (e *ExHttpClient) Serve(ctx context.Context, input *model.Event) (*model.Event, error) {

	resp, err := e.mHttpClient.Get(ctx, input.Metadata, map[string]string{"Content-Type": "application/json"}, e.mUrl)

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return model.GenerateOutputEvent(input.Metadata, e.ContractId(), "OK", 200, "application/json", data), nil
}

func init() {
	registry.GlobalRegistry().AddCapability(&ExHttpClient{})
}
