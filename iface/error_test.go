package iface

import (
	"errors"
	"testing"
)

type internalError struct {
	message     string
	errorCode   uint32
	errorPrefix string
}

func (i *internalError) ErrorPrefix() string {
	return i.errorPrefix
}

func (i *internalError) Error() string {
	return i.message
}

func (i *internalError) ErrorCode() uint32 {
	return i.errorCode
}

func (i *internalError) Message() string {
	return i.message
}

func newInternalError(errorCode uint32, message string) IError {
	return &internalError{
		message:     message,
		errorCode:   errorCode,
		errorPrefix: "TEST_ERROR",
	}
}

var errZero = errors.New("error zero")
var errOne = newInternalError(1, "error one")
var errTwo = newInternalError(2, "error two")

func TestGetErrorInterface(t *testing.T) {
	v := GetErrorInterface(errOne)

	if v == nil {
		t.Error("error should not be nil")
	}

	if v.ErrorCode() != 1 {
		t.Error("error code should be 1")
	}

	v2 := GetErrorInterface(errZero)

	if v2 != nil {
		t.Error("v2 should be nil")
	}
}

func TestIsErrorInterfaceOk(t *testing.T) {
	if IsErrorInterfaceOk(nil) != false {
		t.Error("output should be false")
	}

	if IsErrorInterfaceOk(errOne) != true {
		t.Error("output should be true")
	}
}

func TestStatusCode(t *testing.T) {
	if StatusCode(errOne) != "TEST_ERROR_1" {
		t.Error("status code should be TEST_ERROR_1")
	}
}

func TestIs(t *testing.T) {
	if Is(errOne, errOne) != true {
		t.Error("output should be true")
	}

	if Is(errOne, errTwo) != false {
		t.Error("output should be false")
	}
}
