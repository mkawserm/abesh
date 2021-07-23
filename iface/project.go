package iface

type IProject interface {
	Name() string
	Version() string
	Authors() []string

	ShortDescription() string
	LongDescription() string
}
