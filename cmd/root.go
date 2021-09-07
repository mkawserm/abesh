package cmd

import (
	"fmt"
	"github.com/mkawserm/abesh/conf"
	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/logger"
	"github.com/mkawserm/abesh/registry"
	"github.com/spf13/cobra"
	"os"
)

var rootCMD = &cobra.Command{
	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
		os.Exit(0)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
}

// AddCommand slices to root cmd
func AddCommand(cmdList ...*cobra.Command) {
	rootCMD.AddCommand(cmdList...)
}

func initConfig() {
	if conf.EnvironmentConfigIns().CMDLogEnabled {
		logger.L(constant.Name).Debug("initializing config")
	}

	registry.GlobalRegistry()

	if conf.EnvironmentConfigIns().CMDLogEnabled {
		logger.L(constant.Name).Debug("global registry initialized")
	}
}

func Execute() {
	DefaultCMDHandler(rootCMD)
	rootCMD.AddCommand(abeshCMD)

	if err := rootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
