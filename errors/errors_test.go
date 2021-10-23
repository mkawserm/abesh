package errors

import (
	"errors"
	"testing"
)

var ErrTest = New(1, "T", "failed", map[string]string{}, false)

func TestError_ProtoErrorWithStack(t *testing.T) {
	t.Logf("%+v\n", ErrTest.StackString())
	t.Logf("%+v\n", NewInternalWithCause(errors.New("test"), "just test", nil, "missing").StackString())
}
