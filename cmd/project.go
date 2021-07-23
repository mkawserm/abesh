package cmd

import "github.com/mkawserm/abesh/constant"

type Project struct {

}

func (p *Project) Name() string {
	return constant.Name
}

func (p *Project) Version() string {
	return constant.Version
}

func (p *Project) Authors() []string {
	return constant.Authors
}

func (p *Project) ShortDescription() string {
	return constant.ShortDescription
}

func (p *Project) LongDescription() string {
	return constant.LongDescription
}
