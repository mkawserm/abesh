package utility

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/mkawserm/abesh/model"
)

var manifestTo = `
version: "1"

capabilities:
  - contract_id: "abesh:health"
  - contract_id: "abesh:httpserver"
    values:
      host: "0.0.0.0"
      port: "8080"
      default_request_timeout: "5s"

  - contract_id: "abesh:ex_echo"
    values:
      key1: "test1"
      key2: "test2"

  - contract_id: "abesh:ex_echo"
    new_contract_id: "abesh:ex_echo:1"
    values:
      key1: "test4"
      key2: "test5"

  - contract_id: "abesh:ex_echo"
    new_contract_id: "abesh:ex_echo:2"
    values:
      key1: "test6"
      key2: "test7"

  - contract_id: "abesh:httpserver"
    new_contract_id: "abesh:httpserver:1"
    values:
      host: "0.0.0.0"
      port: "9090"
      default_request_timeout: "5s"

triggers:
  - trigger: "abesh:httpserver"
    trigger_values:
      method: "GET"
      path: "/default1"
    service: "abesh:ex_echo:1"

  - trigger: "abesh:httpserver"
    trigger_values:
      method: "GET"
      path: "/default2"
    service: "abesh:ex_echo:2"

  - trigger: "abesh:httpserver:1"
    trigger_values:
      method: "GET"
      path: "/default1"
    service: "abesh:ex_echo"

start:
  - "abesh:httpserver"
  - "abesh:httpserver:1"
`

var manifestFrom = `
version: "1"

capabilities:
  - contract_id: "abesh:httpserver"
    values:
      host: "0.0.0.0"
      port: "8888"
      default_request_timeout: "10s"

  - contract_id: "abesh:ex_echo"
    new_contract_id: "abesh:ex_echo:1"
    values:
      key1: "test14"
      key2: "test15"

  - contract_id: "abesh:ex_echo"
    new_contract_id: "abesh:ex_echo:2"
    values:
      key1: "test16"
      key2: "test17"

  - contract_id: "abesh:httpserver"
    new_contract_id: "abesh:httpserver:1"
    values:
      port: "9999"
      default_request_timeout: "15s"
`

var manifest = `
version: "1"

capabilities:
  - contract_id: "abesh:health"
  - contract_id: "abesh:httpserver"
    values:
      host: "0.0.0.0"
      port: "8888"
      default_request_timeout: "10s"

  - contract_id: "abesh:ex_echo"
    values:
      key1: "test1"
      key2: "test2"

  - contract_id: "abesh:ex_echo"
    new_contract_id: "abesh:ex_echo:1"
    values:
      key1: "test14"
      key2: "test15"

  - contract_id: "abesh:ex_echo"
    new_contract_id: "abesh:ex_echo:2"
    values:
      key1: "test16"
      key2: "test17"

  - contract_id: "abesh:httpserver"
    new_contract_id: "abesh:httpserver:1"
    values:
      host: "0.0.0.0"
      port: "9999"
      default_request_timeout: "15s"

triggers:
  - trigger: "abesh:httpserver"
    trigger_values:
      method: "GET"
      path: "/default1"
    service: "abesh:ex_echo:1"

  - trigger: "abesh:httpserver"
    trigger_values:
      method: "GET"
      path: "/default2"
    service: "abesh:ex_echo:2"

  - trigger: "abesh:httpserver:1"
    trigger_values:
      method: "GET"
      path: "/default1"
    service: "abesh:ex_echo"

start:
  - "abesh:httpserver"
  - "abesh:httpserver:1"
`

func TestMergeManifest(t *testing.T) {
	manifestModelFrom := new(model.Manifest)
	manifestModelTo := new(model.Manifest)
	manifestModel := new(model.Manifest)
	var err = yaml.Unmarshal([]byte(manifestFrom), manifestModelFrom)
	if err != nil {
		t.Error(err)
	}

	err = yaml.Unmarshal([]byte(manifestTo), manifestModelTo)
	if err != nil {
		t.Error(err)
	}

	var manifestResult = MergeManifest(manifestModelTo, manifestModelFrom)
	err = yaml.Unmarshal([]byte(manifest), manifestModel)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(manifestResult, manifestModel) {
		t.Error("manifest must be equal")
	}
}
