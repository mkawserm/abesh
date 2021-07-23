package cmd

import (
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/platform"
	"github.com/spf13/cobra"
)

type OnRootInitHandler func(cmd *cobra.Command)

var DefaultProject iface.IProject
var DefaultPlatform iface.IPlatform

var DefaultCMDHandler OnRootInitHandler = func(cmd *cobra.Command) {
	DefaultProject = &Project{}
	DefaultPlatform = &platform.One{}

	cmd.Use = DefaultProject.Name()
	cmd.Short = DefaultProject.ShortDescription()
	cmd.Long = DefaultProject.LongDescription()
}
