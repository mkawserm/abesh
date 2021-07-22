package iface

type ICapabilityRegistry interface {
	KVStore(contactId string) IKVStore
	Interface(contractId string) interface{}
}
