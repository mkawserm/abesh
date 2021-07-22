package httpserver

import "github.com/mkawserm/abesh/constant"

type HTTPServer struct {
	mValues map[string]interface{}
}

func (h *HTTPServer) Source() string {
	return "github.com/mkawserm/abesh/capability/httpserver"
}

func (h *HTTPServer) Runtime() string {
	return string(constant.Native)
}

func (h *HTTPServer) Category() string {
	return string(constant.Trigger)
}

func (h *HTTPServer) ContractId() string {
	return "abesh:httpserver"
}

func (h *HTTPServer) Values() map[string]interface{} {
	return h.mValues
}

func (h *HTTPServer) SetValues(values map[string]interface{}) {
	h.mValues = values
}
