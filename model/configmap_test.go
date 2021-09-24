package model

import (
	"reflect"
	"testing"
	"time"
)

func TestConfigMap_String(t *testing.T) {
	c := ConfigMap{"1": "one"}

	if c.String("1", "default") != "one" {
		t.Error("1 should be one")
	}

	if c.String("2", "default") != "default" {
		t.Error("2 should be default")
	}
}

func TestConfigMap_StringList(t *testing.T) {
	c := ConfigMap{"list": "1,2,3"}

	if !reflect.DeepEqual(c.StringList("list", ",", []string{}), []string{"1", "2", "3"}) {
		t.Error("list should match")
	}

	if !reflect.DeepEqual(c.StringList("list1", ",", []string{"1", "2", "3"}), []string{"1", "2", "3"}) {
		t.Error("list should match")
	}
}

func TestConfigMap_StringMap(t *testing.T) {
	c := ConfigMap{"map": "1=1;2=2;3=3"}
	if !reflect.DeepEqual(c.StringMap("map", ConfigMap{}), ConfigMap{"1": "1", "2": "2", "3": "3"}) {
		t.Error("map should match")
	}

	if !reflect.DeepEqual(c.StringMap("map1", ConfigMap{"1": "1", "2": "2", "3": "3"}), ConfigMap{"1": "1", "2": "2", "3": "3"}) {
		t.Error("map should match")
	}
}

func TestConfigMap_Bytes(t *testing.T) {
	c := ConfigMap{"a": "1,2,3"}

	if !reflect.DeepEqual(c.Bytes("a", nil), []byte("1,2,3")) {
		t.Error("a should match")
	}

	if !reflect.DeepEqual(c.Bytes("b", []byte("1,2,3")), []byte("1,2,3")) {
		t.Error("b should match")
	}
}

func TestConfigMap_Float32(t *testing.T) {
	v := ConfigMap{"n": "100.10", "n2": "one"}

	if v.Float32("n", 0) != 100.10 {
		t.Error("n should be 100.10")
	}

	if v.Float32("n2", 1) != 1 {
		t.Error("output should be default 1")
	}

	if v.Float32("n3", 1) != 1 {
		t.Error("output should be default 1")
	}
}

func TestConfigMap_Float64(t *testing.T) {
	v := ConfigMap{"n": "100.10", "n2": "one"}

	if v.Float64("n", 0) != 100.10 {
		t.Error("n should be 100.10")
	}

	if v.Float64("n2", 1) != 1 {
		t.Error("output should be default 1")
	}

	if v.Float64("n3", 1) != 1 {
		t.Error("output should be default 1")
	}
}

func TestConfigMap_Bool(t *testing.T) {
	v := ConfigMap{"b": "false", "a": "true", "d": "one"}

	if v.Bool("b", true) != false {
		t.Error("b should be false")
	}

	if v.Bool("a", false) != true {
		t.Error("a should be true")
	}

	if v.Bool("c", false) != false {
		t.Error("c should be false")
	}

	if v.Bool("d", false) != false {
		t.Error("d should be false")
	}
}

func TestConfigMap_Int(t *testing.T) {
	v := ConfigMap{"n": "100", "n2": "one"}

	if v.Int("n", 0) != 100 {
		t.Error("n should be 100")
	}

	if v.Int("n2", 1) != 1 {
		t.Error("output should be default 1")
	}

	if v.Int("n3", 1) != 1 {
		t.Error("output should be default 1")
	}
}

func TestConfigMap_Uint(t *testing.T) {
	v := ConfigMap{"n": "100", "n2": "one"}

	if v.Uint("n", 0) != 100 {
		t.Error("n should be 100")
	}

	if v.Uint("n2", 1) != 1 {
		t.Error("output should be default 1")
	}

	if v.Uint("n3", 1) != 1 {
		t.Error("output should be default 1")
	}
}

func TestConfigMap_Int16(t *testing.T) {
	v := ConfigMap{"n": "100", "n2": "one"}

	if v.Int16("n", 0) != 100 {
		t.Error("n should be 100")
	}

	if v.Int16("n2", 1) != 1 {
		t.Error("output should be default 1")
	}

	if v.Int16("n3", 1) != 1 {
		t.Error("output should be default 1")
	}
}

func TestConfigMap_Uint16(t *testing.T) {
	v := ConfigMap{"n": "100", "n2": "one"}

	if v.Uint16("n", 0) != 100 {
		t.Error("n should be 100")
	}

	if v.Uint16("n2", 1) != 1 {
		t.Error("output should be default 1")
	}

	if v.Uint16("n3", 1) != 1 {
		t.Error("output should be default 1")
	}
}

func TestConfigMap_Int32(t *testing.T) {
	v := ConfigMap{"n": "100", "n2": "one"}

	if v.Int32("n", 0) != 100 {
		t.Error("n should be 100")
	}

	if v.Int32("n2", 1) != 1 {
		t.Error("output should be default 1")
	}

	if v.Int32("n3", 1) != 1 {
		t.Error("output should be default 1")
	}
}

func TestConfigMap_Uint32(t *testing.T) {
	v := ConfigMap{"n": "100", "n2": "one"}

	if v.Uint32("n", 0) != 100 {
		t.Error("n should be 100")
	}

	if v.Uint32("n2", 1) != 1 {
		t.Error("output should be default 1")
	}

	if v.Uint32("n3", 1) != 1 {
		t.Error("output should be default 1")
	}
}

func TestConfigMap_Int64(t *testing.T) {
	v := ConfigMap{"n": "100", "n2": "one"}

	if v.Int64("n", 0) != 100 {
		t.Error("n should be 100")
	}

	if v.Int64("n2", 1) != 1 {
		t.Error("output should be default 1")
	}

	if v.Int64("n3", 1) != 1 {
		t.Error("output should be default 1")
	}
}

func TestConfigMap_Uint64(t *testing.T) {
	v := ConfigMap{"n": "100", "n2": "one"}

	if v.Uint64("n", 0) != 100 {
		t.Error("n should be 100")
	}

	if v.Uint64("n2", 1) != 1 {
		t.Error("output should be default 1")
	}

	if v.Uint64("n3", 1) != 1 {
		t.Error("output should be default 1")
	}
}

func TestConfigMap_Duration(t *testing.T) {
	c := ConfigMap{"a": "1ms", "b": "one"}

	if c.Duration("a", 2*time.Millisecond) != time.Millisecond {
		t.Error("a should be 1ms")
	}

	if c.Duration("b", time.Millisecond) != time.Millisecond {
		t.Error("b should be 1ms")
	}

	if c.Duration("c", 2*time.Millisecond) != 2*time.Millisecond {
		t.Error("c should be 2ms")
	}
}

func TestConfigMap_Time(t *testing.T) {
	c := ConfigMap{
		"time":  "2002-10-02T15:00:00.05Z",
		"time2": "time2",
	}
	_t, _ := time.Parse(time.RFC3339, "2002-10-02T15:00:00.05Z")

	if c.Time("time", time.Now()) != _t {
		t.Error("time should match")
	}

	if c.Time("time1", _t) != _t {
		t.Error("time should match")
	}

	if c.Time("time2", _t) != _t {
		t.Error("time should match")
	}
}
