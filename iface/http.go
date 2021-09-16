package iface

import (
	"context"
	"github.com/mkawserm/abesh/model"
	"io"
	"net/http"
)

type IHttpClient interface {
	Get(ctx context.Context, metadata *model.Metadata, headers map[string]string, url string) (*http.Response, error)
	Head(ctx context.Context, metadata *model.Metadata, headers map[string]string, url string) (*http.Response, error)
	Delete(ctx context.Context, metadata *model.Metadata, headers map[string]string, url string) (*http.Response, error)
	Post(ctx context.Context, metadata *model.Metadata, headers map[string]string, url string, body io.Reader) (*http.Response, error)
	Do(ctx context.Context, method string, metadata *model.Metadata, headers map[string]string, url string, body io.Reader) (*http.Response, error)
}
