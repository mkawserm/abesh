package stack

import "fmt"

// String transforms stack trace to string
func String(s Stack) string {
	stackStr := ""
	for _, frame := range s {
		stackStr = fmt.Sprintf("%s\n  %s:%d in %s", stackStr, frame.Filename, frame.Line, frame.Method)
	}
	return stackStr
}
