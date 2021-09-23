package httpclient

import (
	"context"
	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
	"github.com/mkawserm/abesh/utility"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type HTTPClient struct {
	mValues model.ConfigMap

	mDialerTimeout       time.Duration
	mTLSHandshakeTimeout time.Duration
	mRequestTimeout      time.Duration

	mDisableKeepAlive   bool
	mDisableCompression bool

	mMaxIdleConnections        int
	mMaxIdleConnectionsPerHost int
	mMaxConnectionsPerHost     int

	mIdleConnectionTimeout time.Duration
	mResponseHeaderTimeout time.Duration
	mExpectContinueTimeout time.Duration

	mMaxResponseHeaderBytes int64
	mWriteBufferSize        int
	mReadBufferSize         int

	mHttpClient *http.Client
}

func (h *HTTPClient) Name() string {
	return "abesh_httpclient"
}

func (h *HTTPClient) Version() string {
	return constant.Version
}

func (h *HTTPClient) Category() string {
	return string(constant.CategoryNetwork)
}

func (h *HTTPClient) ContractId() string {
	return "abesh:httpclient"
}

func (h *HTTPClient) GetConfigMap() model.ConfigMap {
	return h.mValues
}

func (h *HTTPClient) SetConfigMap(values model.ConfigMap) error {
	h.mValues = values

	h.mDialerTimeout = h.mValues.Duration("dialer_timeout", 5*time.Second)
	h.mTLSHandshakeTimeout = h.mValues.Duration("tls_handshake_timeout", 5*time.Second)
	h.mRequestTimeout = h.mValues.Duration("request_timeout", 10*time.Second)
	h.mDisableKeepAlive = h.mValues.Bool("disable_keep_alive", false)
	h.mDisableCompression = h.mValues.Bool("disable_compression", false)
	h.mMaxIdleConnections = h.mValues.Int("max_idle_connections", 100)
	h.mMaxIdleConnectionsPerHost = h.mValues.Int("max_idle_connections_per_host", 10)
	h.mMaxConnectionsPerHost = h.mValues.Int("max_connections_per_host", 1000)
	h.mIdleConnectionTimeout = h.mValues.Duration("idle_connection_timeout", 5*time.Second)
	h.mResponseHeaderTimeout = h.mValues.Duration("response_header_timeout", 0*time.Second)
	h.mExpectContinueTimeout = h.mValues.Duration("expect_continue_timeout", 0*time.Second)
	h.mMaxResponseHeaderBytes = h.mValues.Int64("max_response_header_bytes", 0)
	h.mWriteBufferSize = h.mValues.Int("write_buffer_size", 4096)
	h.mReadBufferSize = h.mValues.Int("read_buffer_size", 4096)

	return nil
}

func (h *HTTPClient) Setup() error {
	dialer := &net.Dialer{
		Timeout: h.mDialerTimeout,
	}

	transport := &http.Transport{
		DialContext: dialer.DialContext,
		//DialTLSContext:      dialer.DialContext,
		TLSHandshakeTimeout: h.mTLSHandshakeTimeout,

		DisableKeepAlives:   h.mDisableKeepAlive,
		DisableCompression:  h.mDisableCompression,
		MaxIdleConns:        h.mMaxIdleConnections,
		MaxIdleConnsPerHost: h.mMaxIdleConnectionsPerHost,
		MaxConnsPerHost:     h.mMaxConnectionsPerHost,

		ResponseHeaderTimeout: h.mResponseHeaderTimeout,
		ExpectContinueTimeout: h.mExpectContinueTimeout,

		MaxResponseHeaderBytes: h.mMaxResponseHeaderBytes,
		WriteBufferSize:        h.mWriteBufferSize,
		ReadBufferSize:         h.mReadBufferSize,
	}

	h.mHttpClient = &http.Client{
		Transport: transport,
		Timeout:   h.mRequestTimeout,
	}

	return nil
}

func (h *HTTPClient) New() iface.ICapability {
	return &HTTPClient{}
}

func (h *HTTPClient) Get(ctx context.Context, metadata *model.Metadata, headers map[string]string, url string) (*http.Response, error) {
	return h.Do(ctx, "GET", metadata, headers, url, nil)
}

func (h *HTTPClient) Head(ctx context.Context, metadata *model.Metadata, headers map[string]string, url string) (*http.Response, error) {
	return h.Do(ctx, "HEAD", metadata, headers, url, nil)
}

func (h *HTTPClient) Delete(ctx context.Context, metadata *model.Metadata, headers map[string]string, url string) (*http.Response, error) {
	return h.Do(ctx, "DELETE", metadata, headers, url, nil)
}

func (h *HTTPClient) Post(ctx context.Context, metadata *model.Metadata, headers map[string]string, url string, body io.Reader) (*http.Response, error) {
	return h.Do(ctx, "POST", metadata, headers, url, body)
}

func (h *HTTPClient) Do(ctx context.Context, method string, metadata *model.Metadata, headers map[string]string, url string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequestWithContext(ctx, strings.ToUpper(method), url, body)

	mergedHeaders := headers
	if metadata != nil {
		mergedHeaders = utility.Merge(metadata.Headers, headers)
	}

	if mergedHeaders != nil {
		for k, v := range mergedHeaders {
			r.Header.Set(k, v)
		}
	}

	if err != nil {
		return nil, err
	}

	return h.mHttpClient.Do(r)
}

func init() {
	registry.GlobalRegistry().AddCapability(&HTTPClient{})
}

// GetHttpClient returns http client capability
func GetHttpClient(registry iface.ICapabilityRegistry) iface.IHttpClient {
	c := registry.Capability("abesh:httpclient")

	if c == nil {
		return nil
	}

	c2, ok := c.(*HTTPClient)

	if !ok {
		return nil
	}

	return c2
}
