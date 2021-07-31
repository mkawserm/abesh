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
	"strconv"
	"strings"
	"time"
)

type HTTPClient struct {
	mValues map[string]string

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
	return "golang_net_http_client"
}

func (h *HTTPClient) Version() string {
	return "0.0.1"
}

func (h *HTTPClient) Category() string {
	return string(constant.CategoryNetwork)
}

func (h *HTTPClient) ContractId() string {
	return "abesh:httpclient"
}

func (h *HTTPClient) Values() map[string]string {
	return h.mValues
}

func (h *HTTPClient) SetValues(values map[string]string) error {
	h.mValues = values

	var t time.Duration
	var i int
	var i64 int64
	var b bool
	var err error

	// dialer timeout
	t, err = time.ParseDuration(utility.GetValue(h.mValues, "dialer_timeout", "5s"))
	if err != nil {
		h.mDialerTimeout = 5 * time.Second
	} else {
		h.mDialerTimeout = t
	}

	//tls handshake timeout
	t, err = time.ParseDuration(utility.GetValue(h.mValues, "tls_handshake_timeout", "5s"))
	if err != nil {
		h.mTLSHandshakeTimeout = 5 * time.Second
	} else {
		h.mTLSHandshakeTimeout = t
	}

	//request timeout
	t, err = time.ParseDuration(utility.GetValue(h.mValues, "request_timeout", "10s"))
	if err != nil {
		h.mTLSHandshakeTimeout = 10 * time.Second
	} else {
		h.mTLSHandshakeTimeout = t
	}

	// disable keep alive
	b, err = strconv.ParseBool(utility.GetValue(h.mValues, "disable_keep_alive", "false"))
	if err != nil {
		h.mDisableKeepAlive = false
	} else {
		h.mDisableKeepAlive = b
	}

	// disable compression
	b, err = strconv.ParseBool(utility.GetValue(h.mValues, "disable_compression", "false"))
	if err != nil {
		h.mDisableCompression = false
	} else {
		h.mDisableCompression = b
	}

	// max idle connections
	i, err = strconv.Atoi(utility.GetValue(h.mValues, "max_idle_connections", "100"))
	if err != nil {
		h.mMaxIdleConnections = 100
	} else {
		h.mMaxIdleConnections = i
	}

	// max idle connections per host
	i, err = strconv.Atoi(utility.GetValue(h.mValues, "max_idle_connections_per_host", "10"))
	if err != nil {
		h.mMaxIdleConnectionsPerHost = 10
	} else {
		h.mMaxIdleConnectionsPerHost = i
	}

	// max connections per host
	i, err = strconv.Atoi(utility.GetValue(h.mValues, "max_connections_per_host", "1000"))
	if err != nil {
		h.mMaxConnectionsPerHost = 1000
	} else {
		h.mMaxConnectionsPerHost = i
	}

	// idle connection timeout
	t, err = time.ParseDuration(utility.GetValue(h.mValues, "idle_connection_timeout", "5s"))
	if err != nil {
		h.mIdleConnectionTimeout = 5 * time.Second
	} else {
		h.mIdleConnectionTimeout = t
	}

	// response header timeout
	t, err = time.ParseDuration(utility.GetValue(h.mValues, "response_header_timeout", "0s"))
	if err != nil {
		h.mResponseHeaderTimeout = 0 * time.Second
	} else {
		h.mResponseHeaderTimeout = t
	}

	// expect continue timeout
	t, err = time.ParseDuration(utility.GetValue(h.mValues, "expect_continue_timeout", "0s"))
	if err != nil {
		h.mExpectContinueTimeout = 0 * time.Second
	} else {
		h.mExpectContinueTimeout = t
	}

	// max response header bytes
	i64, err = strconv.ParseInt(utility.GetValue(h.mValues, "max_response_header_bytes", "0"), 10, 64)
	if err != nil {
		h.mMaxResponseHeaderBytes = 0
	} else {
		h.mMaxResponseHeaderBytes = i64
	}

	// write buffer size
	i, err = strconv.Atoi(utility.GetValue(h.mValues, "write_buffer_size", "4096"))
	if err != nil {
		h.mWriteBufferSize = 4096
	} else {
		h.mWriteBufferSize = i
	}

	// read buffer size
	i, err = strconv.Atoi(utility.GetValue(h.mValues, "read_buffer_size", "4096"))
	if err != nil {
		h.mReadBufferSize = 4096
	} else {
		h.mReadBufferSize = i
	}

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
