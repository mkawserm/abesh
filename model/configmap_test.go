package model

import "testing"

func TestConfigMap_Bool(t *testing.T) {
	v := ConfigMap{"b": "false", "a": "true"}

	if v.Bool("b", true) != false {
		t.Error("b should be false")
	}

	if v.Bool("a", false) != true {
		t.Error("a should be true")
	}
}

func TestConfigMap_Int(t *testing.T) {
	v := ConfigMap{"n": "100"}

	if v.Int("n", 0) != 100 {
		t.Error("n should be 100")
	}
}
