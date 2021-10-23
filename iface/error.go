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

	GetMessage() string
	GetCode() uint32
	GetPrefix() string
	GetParams() map[string]string
	IsRetryable() bool
}

// StatusCode2 generate error status code from IError2
func StatusCode2(iError IError2) string {
	return fmt.Sprintf("%s_%d", iError.GetPrefix(), iError.GetCode())
}

// GetError2 type cast to IError2
func GetError2(anyError interface{}) IError2 {
	iError, ok := anyError.(IError2)
	if ok {
		return iError
	}
	return nil
}

// Is2 checks error with target error
func Is2(err, target IError2) bool {
	if err.GetCode() == target.GetCode() && err.GetPrefix() == target.GetPrefix() {
		return true
	}

	return false
}
