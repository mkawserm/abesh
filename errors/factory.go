package errors

import (
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/stack"
	"strings"
)

func NewInternalErrorWithCause(err error, message string, params map[string]string, subCode string) *Error {
	newErr := errorFactory(444, errCode(ErrInternalService, subCode), message, params)
	newErr.cause = err

	switch v := err.(type) {
	case *Error:
		if v.errorModel.Retryable != nil {
			newErr.errorModel.Retryable = v.errorModel.Retryable
		}
	}

	return newErr
}

func errorFactory(code uint32, prefix string, message string, params map[string]string) *Error {
	e := &Error{
		errorModel: &model.Error{
			Status: &model.Status{
				Code:    code,
				Prefix:  ErrUnknown,
				Message: message,
				Params:  params,
			},
			Retryable: &model.BoolValue{
				Value: false,
			},
		},
	}

	if len(prefix) > 0 {
		e.errorModel.Status.Prefix = prefix

		e.errorModel.Retryable.Value = false
		for _, c := range retryableCodes {
			if PrefixMatches(e, c) {
				e.errorModel.Retryable.Value = true
			}
		}
	}

	if params != nil {
		e.errorModel.Status.Params = params
	}

	e.stackFrames = stack.BuildStack(3)

	return e
}

// InternalService creates a new error to represent an internal service error.
// Only use internal service error if we know very little about the error. Most
// internal service errors will come from `Wrap`ing a vanilla `error` interface
func InternalService(prefix, message string, params map[string]string) *Error {
	return errorFactory(424, errCode(ErrInternalService, prefix), message, params)
}

// BadRequest creates a new error to represent an error caused by the client sending
// an invalid request. This is non-retryable unless the request is modified.
func BadRequest(prefix, message string, params map[string]string) *Error {
	return errorFactory(400, errCode(ErrBadRequest, prefix), message, params)
}

// BadResponse creates a new error representing a failure to response with a valid response
// Examples of this would be a handler returning an invalid message format
func BadResponse(prefix, message string, params map[string]string) *Error {
	return errorFactory(400, errCode(ErrBadResponse, prefix), message, params)
}

// Timeout creates a new error representing a timeout from client to server
func Timeout(prefix, message string, params map[string]string) *Error {
	return errorFactory(408, errCode(ErrTimeout, prefix), message, params)
}

// NotFound creates a new error representing a resource that cannot be found. In some
// cases this is not an error, and would be better represented by a zero length slice of elements
func NotFound(prefix, message string, params map[string]string) *Error {
	return errorFactory(404, errCode(ErrNotFound, prefix), message, params)
}

// Forbidden creates a new error representing a resource that cannot be accessed with
// the current authorisation credentials. The user may need authorising, or if authorised,
// may not be permitted to perform this action
func Forbidden(prefix, message string, params map[string]string) *Error {
	return errorFactory(403, errCode(ErrForbidden, prefix), message, params)
}

// Unauthorized creates a new error indicating that authentication is required,
// but has either failed or not been provided.
func Unauthorized(prefix, message string, params map[string]string) *Error {
	return errorFactory(401, errCode(ErrUnauthorized, prefix), message, params)
}

// PreconditionFailed creates a new error indicating that one or more conditions
// given in the request evaluated to false when tested on the server.
func PreconditionFailed(prefix, message string, params map[string]string) *Error {
	return errorFactory(412, errCode(ErrPreconditionFailed, prefix), message, params)
}

// RateLimited creates a new error indicating that the request has been rate-limited,
// and that the caller should back-off.
func RateLimited(prefix, message string, params map[string]string) *Error {
	return errorFactory(429, errCode(ErrRateLimited, prefix), message, params)
}

func errCode(prefix, sub string) string {
	if sub == "" {
		return prefix
	}
	if prefix == "" {
		return sub
	}
	return strings.Join([]string{prefix, sub}, ".")
}

func NewWithAllInfo(code uint32, prefix string, message string, params map[string]string, retryable bool) *Error {
	e := &Error{
		errorModel: &model.Error{
			Status: &model.Status{
				Code:    code,
				Prefix:  prefix,
				Message: message,
				Params:  params,
			},
			Retryable: &model.BoolValue{
				Value: retryable,
			},
		},
	}

	e.stackFrames = stack.BuildStack(2)

	return e
}

func NewFromError(source *Error) *Error {
	e := &Error{
		errorModel: &model.Error{
			Status: &model.Status{
				Code:    source.GetCode(),
				Prefix:  source.GetPrefix(),
				Message: source.GetMessage(),
				Params:  source.GetParams(),
			},
			Retryable: &model.BoolValue{
				Value: source.IsRetryable(),
			},
		},
	}

	e.stackFrames = stack.BuildStack(2)
	return e
}

// New creates a new error for you. Use this if you want to pass along a custom error code.
// Otherwise, use the handy shorthand factories below
func New(code uint32, prefix string, message string, params map[string]string) *Error {
	return errorFactory(code, prefix, message, params)
}
