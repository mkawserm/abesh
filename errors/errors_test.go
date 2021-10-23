package errors

import (
	"errors"
	"testing"
)

var ErrTest = NewError(1, "T", "failed", map[string]string{}, true)

func TestError_ProtoErrorWithStack(t *testing.T) {
	t.Logf("%+v\n", ErrTest.StackString())
	t.Logf("%+v\n", NewInternalWithCause(errors.New("test"), "just test", nil, "missing").StackString())

	t.Logf("%+v\n", NewErrorFromSource(ErrTest).StackString())
}

func TestError_IsRetryable(t *testing.T) {
	if !IsRetryable(ErrTest) {
		t.Error("error should be retryable")
	}
}

func Test_IsPrefixMatches(t *testing.T) {
	if !IsPrefixMatches(ErrTest, "T") {
		t.Error("prefix should match with T")
	}
}
