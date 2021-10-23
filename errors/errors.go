// idea taken from https://github.com/monzo/terrors/blob/master/errors.go

package errors

import (
	"fmt"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/stack"
	"strings"
)

// Generic error codes. Each of these has their own constructor for convenience.
// You can use any string as a code, just use the `New` method.
const (
	ErrBadRequest         = "bad_request"
	ErrBadResponse        = "bad_response"
	ErrForbidden          = "forbidden"
	ErrInternalService    = "internal_service"
	ErrNotFound           = "not_found"
	ErrPreconditionFailed = "precondition_failed"
	ErrTimeout            = "timeout"
	ErrUnauthorized       = "unauthorized"
	ErrUnknown            = "unknown"
	ErrRateLimited        = "rate_limited"
)

var retryableCodes = []string{
	ErrInternalService,
	ErrTimeout,
	ErrUnknown,
	ErrRateLimited,
}

type Error struct {
	errorModel  *model.Error
	stackFrames stack.Stack
	cause       error
}

func (x *Error) Error() string {
	if x.cause == nil {
		return x.legacyErrString()
	}

	var next error = x
	output := strings.Builder{}
	output.WriteString(x.GetPrefix())
	for next != nil {
		output.WriteString(": ")
		switch typed := next.(type) {
		case *Error:
			output.WriteString(typed.GetMessage())
			next = typed.cause
		case error:
			output.WriteString(typed.Error())
			next = nil
		}
	}
	return output.String()
}

// Unwrap returns the cause of the error. It may be nil.
func (x *Error) Unwrap() error {
	return x.cause
}

func (x *Error) legacyErrString() string {
	if x == nil {
		return ""
	}
	if x.errorModel == nil {
		return ""
	}

	if x.errorModel.Status == nil {
		return ""
	}

	if x.errorModel.Status.Message == "" {
		return x.errorModel.Status.Prefix
	}

	if x.errorModel.Status.Prefix == "" {
		return x.errorModel.Status.Message
	}

	return fmt.Sprintf("%s: %s", x.errorModel.Status.Prefix, x.errorModel.Status.Message)
}

// StackTrace returns a slice of program counters taken from the stack frames.
// This adapts the terrors package to allow stacks to be reported to Sentry correctly.
func (x *Error) StackTrace() []uintptr {
	out := make([]uintptr, len(x.stackFrames))
	for i := 0; i < len(x.stackFrames); i++ {
		out[i] = x.stackFrames[i].PC
	}
	return out
}

// StackString formats the stack as a beautiful string with newlines
func (x *Error) StackString() string {
	stackStr := ""
	for _, frame := range x.stackFrames {
		stackStr = fmt.Sprintf("%s\n  %s:%d in %s", stackStr, frame.Filename, frame.Line, frame.Method)
	}
	return stackStr
}

// VerboseString returns the error message, stack trace and params
func (x *Error) VerboseString() string {
	return fmt.Sprintf("%s\nParams: %+v\n%s", x.Error(), x.GetParams(), x.StackString())
}

func (x *Error) PrefixMatches(prefixParts ...string) bool {
	prefix := strings.Join(prefixParts, ".")

	return strings.HasPrefix(x.GetPrefix(), prefix)
}

func PrefixMatches(err error, prefixParts ...string) bool {
	if terr, ok := Wrap(err, nil).(*Error); ok {
		return terr.PrefixMatches(prefixParts...)
	}

	return false
}

// Retryable determines whether the error was caused by an action which can be retried.
func (x *Error) Retryable() bool {
	if x.errorModel.Retryable != nil {
		return x.errorModel.Retryable.Value
	}

	for _, c := range retryableCodes {
		if PrefixMatches(x, c) {
			return true
		}
	}
	return false
}

func (x *Error) LogMetadata() map[string]string {
	return x.GetParams()
}

func (x *Error) GetMessage() string {
	return x.errorModel.Status.GetMessage()
}

func (x *Error) GetPrefix() string {
	return x.errorModel.Status.GetPrefix()
}

func (x *Error) GetCode() uint32 {
	return x.errorModel.Status.GetCode()
}

func (x *Error) GetParams() map[string]string {
	return x.errorModel.Status.GetParams()
}

func (x *Error) IsRetryable() bool {
	return x.errorModel.Retryable.Value
}

func (x *Error) simpleErrString() string {
	if x == nil {
		return ""
	}

	if x.errorModel == nil {
		return ""
	}

	if x.errorModel.Status == nil {
		return ""
	}

	return fmt.Sprintf("%s: %d - %s",
		x.errorModel.Status.Prefix,
		x.errorModel.Status.Code,
		x.errorModel.Status.Message)
}

func (x *Error) ProtoStack() []*model.StackFrame {
	sf := make([]*model.StackFrame, 0, len(x.stackFrames))
	for _, frame := range x.stackFrames {
		sf = append(sf, &model.StackFrame{
			Filename: frame.Filename,
			Line:     int32(frame.Line),
			Method:   frame.Method,
		})
	}
	return sf
}

func (x *Error) ProtoErrorWithStack() *model.ErrorWithStack {
	return &model.ErrorWithStack{
		Error: x.errorModel,
		Stack: x.ProtoStack(),
	}
}

func (x *Error) FromProtoErrorWithStack(input *model.ErrorWithStack) {
	x.errorModel = input.Error
	x.stackFrames = make(stack.Stack, 0, len(input.Stack))
	for _, frame := range input.Stack {
		x.stackFrames = append(x.stackFrames, &stack.Frame{
			Filename: frame.Filename,
			Method:   frame.Method,
			Line:     int(frame.Line),
			PC:       0,
		})
	}
}

func New(code uint32, prefix string, message string, params map[string]string, retryable bool) *Error {
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

	e.stackFrames = stack.BuildStack(3)

	return e
}

// addParams returns a new error with new params merged into the original error's
func addParams(err *Error, params map[string]string) *Error {
	copiedParams := make(map[string]string, len(err.GetParams())+len(params))
	for k, v := range err.GetParams() {
		copiedParams[k] = v
	}
	for k, v := range params {
		copiedParams[k] = v
	}

	e := &Error{errorModel: err.errorModel, stackFrames: err.stackFrames}
	e.errorModel.Status.Params = copiedParams

	return e
}

func IsRetryable(err error) bool {
	if r, ok := Propagate(err).(*Error); ok {
		return r.Retryable()
	}
	return false
}

func Augment(errInput error, context string, params map[string]string) error {
	if errInput == nil {
		return nil
	}
	switch err := errInput.(type) {
	case *Error:
		withMergedParams := addParams(err, params)
		e := &Error{
			errorModel: &model.Error{
				Status: &model.Status{
					Code:    withMergedParams.GetCode(),
					Prefix:  err.GetPrefix(),
					Message: context,
					Params:  withMergedParams.GetParams(),
				},
				Retryable: &model.BoolValue{
					Value: err.IsRetryable(),
				},
			},
			stackFrames: stack.Stack{},
		}
		return e
	default:
		return NewInternalWithCause(err, context, params, "")
	}
}

func Propagate(errInput error) error {
	if errInput == nil {
		return nil
	}
	switch err := errInput.(type) {
	case *Error:
		return err
	default:
		return NewInternalWithCause(err, err.Error(), nil, "")
	}
}

func Is(errInput error, prefix ...string) bool {
	switch err := errInput.(type) {
	case *Error:
		if err.PrefixMatches(prefix...) {
			return true
		}
		next := err.Unwrap()
		if next == nil {
			return false
		}
		return Is(next, prefix...)
	default:
		return false
	}
}
