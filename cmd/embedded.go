package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var embeddedCMD = &cobra.Command{
	Use:   "embedded",
	Short: "Embedded sub command",
	Long:  "Embedded sub command",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		os.Exit(0)
	},
}

func init() {
	embeddedCMD.PersistentFlags().StringSliceVar(&manifestFilePathList, "manifest", []string{}, "Manifest file path list (ex: /home/manifest1.yaml,/home/manifest2.yaml)")
	rootCMD.AddCommand(embeddedCMD)
}
