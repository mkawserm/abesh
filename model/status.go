package model

// NewStatus instantiate new status struct from code, prefix, message and params
func NewStatus(code uint32, prefix string, message string, params map[string]string) *Status {
	return &Status{
		Code:    code,
		Prefix:  prefix,
		Message: message,
		Params:  params,
	}
}
