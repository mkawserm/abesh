package iface

import "fmt"

type IError interface {
	Error() string
	Message() string
	ErrorCode() uint32
	ErrorPrefix() string
}

// StatusCode generate error status code from iError
func StatusCode(iError IError) string {
	return fmt.Sprintf("%s_%d", iError.ErrorPrefix(), iError.ErrorCode())
}

// GetErrorInterface type cast to IError
func GetErrorInterface(anyError interface{}) IError {
	iError, ok := anyError.(IError)
	if ok {
		return iError
	}
	return nil
}

// IsErrorInterfaceOk checks if the IError is
// not nil
func IsErrorInterfaceOk(iError IError) bool {
	if iError == nil {
		return false
	}

	return true
}

// Is checks error with target error
func Is(err, target IError) bool {
	if err.ErrorCode() == target.ErrorCode() && err.ErrorPrefix() == target.ErrorPrefix() {
		return true
	}

	return false
}
