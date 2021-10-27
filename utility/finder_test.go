package utility

import "testing"

func TestIsIn(t *testing.T) {
	if !IsIn([]string{"1", "100", "20"}, "20") {
		t.Errorf("20 should be in the list")
	}

	if IsIn([]string{"1", "100", "20"}, "200") {
		t.Errorf("200 should not be in the list")
	}
}
