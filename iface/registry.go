package iface

type ICapabilityRegistry interface {
	KVStore(contractId string) IKVStore
	Interface(contractId string) interface{}
}
