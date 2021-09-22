package cmd

import (
	"github.com/spf13/cobra"
)

var runCMD = &cobra.Command{
	Use:   "run",
	Short: "run the platform",
	Long:  "run all platform components",
	Run: func(cmd *cobra.Command, args []string) {
		p := PlatformSetup(manifestFilePath)
		p.Run()
	},
}

func init() {
	runCMD.Flags().StringVar(&manifestFilePath, "manifest", "", "Manifest file path (ex: /home/ubuntu/manifest.yaml)")
	_ = runCMD.MarkFlagRequired("manifest")
	rootCMD.AddCommand(runCMD)
}
