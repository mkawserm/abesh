package utility

import (
	httpclient2 "github.com/mkawserm/abesh/capability/httpclient"
	"github.com/mkawserm/abesh/iface"
)

// GetHttpClient returns http client capability
func GetHttpClient(registry iface.ICapabilityRegistry) *httpclient2.HTTPClient {
	c := registry.Capability("abesh:httpclient")

	if c == nil {
		return nil
	}

	c2, ok := c.(*httpclient2.HTTPClient)

	if !ok {
		return nil
	}

	return c2
}
