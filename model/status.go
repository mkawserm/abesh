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

// CloneStatus clone status to new status with new params
func CloneStatus(status *Status, params map[string]string) *Status {
	return &Status{
		Code:    status.Code,
		Prefix:  status.Prefix,
		Message: status.Message,
		Params:  params,
	}
}

// CloneStatusWithMergedParams clone status to new status with new params merged
func CloneStatusWithMergedParams(status *Status, params map[string]string) *Status {
	p := make(map[string]string)

	for k, v := range status.Params {
		p[k] = v
	}

	for k, v := range params {
		p[k] = v
	}

	return &Status{
		Code:    status.Code,
		Prefix:  status.Prefix,
		Message: status.Message,
		Params:  p,
	}
}
