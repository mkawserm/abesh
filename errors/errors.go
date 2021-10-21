package errors

import (
	"fmt"
	abeshStack "github.com/mkawserm/abesh/stack"
	"google.golang.org/protobuf/proto"
)

func (x *Error) Error() string {
	return fmt.Sprintf("%d - %s", x.ErrorCode, x.ErrorMessage)
}

func (x *Error) Message() string {
	return x.GetErrorMessage()
}

func (x *Error) Params() map[string]string {
	return x.GetErrorParams()
}

func (x *Error) Code() uint32 {
	return x.GetErrorCode()
}

func (x *Error) Prefix() string {
	return x.GetErrorPrefix()
}

func (x *Error) IsRetryable() bool {
	if x.ErrorRetryable == nil {
		return false
	}
	return x.ErrorRetryable.Value
}

func (x *Error) Stack(skip int) abeshStack.Stack {
	return abeshStack.BuildStack(skip)
}

func (x *Error) ProtoStack(skip int) []*StackFrame {
	bs := abeshStack.BuildStack(skip)
	sf := make([]*StackFrame, 0, len(bs))
	for _, frame := range bs {
		sf = append(sf, &StackFrame{
			Filename: frame.Filename,
			Line:     int32(frame.Line),
			Method:   frame.Method,
		})
	}

	return sf
}

func (x *Error) ProtoErrorWithStack() *ErrorWithStack {
	return &ErrorWithStack{
		Error: proto.Clone(x).(*Error),
		Stack: x.ProtoStack(3),
	}
}

func New(code uint32, prefix string, message string, retryable bool, params map[string]string) *Error {
	e := &Error{
		ErrorCode:      code,
		ErrorPrefix:    prefix,
		ErrorMessage:   message,
		ErrorParams:    params,
		ErrorRetryable: &BoolValue{Value: retryable},
	}
	return e
}
