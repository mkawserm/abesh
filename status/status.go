package status

import "github.com/mkawserm/abesh/model"

// New instantiate new status struct from code, prefix, message and params
func New(code uint32, prefix string, message string, params map[string]string) *model.Status {
	return &model.Status{
		Code:    code,
		Prefix:  prefix,
		Message: message,
		Params:  params,
	}
}

// Clone status to new status with new params
func Clone(status *model.Status, params map[string]string) *model.Status {
	return &model.Status{
		Code:    status.Code,
		Prefix:  status.Prefix,
		Message: status.Message,
		Params:  params,
	}
}

// CloneWithMergedParams clone status to new status with new params merged
func CloneWithMergedParams(status *model.Status, params map[string]string) *model.Status {
	p := make(map[string]string)

	for k, v := range status.Params {
		p[k] = v
	}

	for k, v := range params {
		p[k] = v
	}

	return &model.Status{
		Code:    status.Code,
		Prefix:  status.Prefix,
		Message: status.Message,
		Params:  p,
	}
}
