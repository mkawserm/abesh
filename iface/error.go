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

type IError2 interface {
	Error() string

	Message() string
	Code() uint32
	Prefix() string
	Params() map[string]string
	IsRetryable() bool
}

// StatusCode2 generate error status code from IError2
func StatusCode2(iError IError2) string {
	return fmt.Sprintf("%s_%d", iError.Prefix(), iError.Code())
}

// GetError2 type cast to IError2
func GetError2(anyError interface{}) IError2 {
	iError, ok := anyError.(IError2)
	if ok {
		return iError
	}
	return nil
}
