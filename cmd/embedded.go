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
	embeddedCMD.PersistentFlags().StringVar(&manifestFilePath, "manifest", "", "Manifest file path (ex: /home/ubuntu/manifest.yaml)")
	rootCMD.AddCommand(embeddedCMD)
}
