package utility

import "strings"

func AsStringList(text string) []string {
	return strings.Split(text, "\n")
}
