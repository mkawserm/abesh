package errors

import "testing"

var ErrTest = New(1, "T_", "failed", false, map[string]string{})

func TestError_ProtoErrorWithStack(t *testing.T) {
	t.Logf("%+v\n", ErrTest.ProtoErrorWithStack())
}
