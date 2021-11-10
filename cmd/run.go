package cmd

import (
	"github.com/spf13/cobra"
)

var runCMD = &cobra.Command{
	Use:   "run",
	Short: "Run the platform",
	Long:  "Run all platform components",
	Run: func(cmd *cobra.Command, args []string) {
		p := PlatformSetup(manifestFilePath)
		p.Run()
	},
}

func init() {
	runCMD.Flags().StringVar(&manifestFilePath, "manifest", "", "Manifest file path (ex: /home/ubuntu/data.txt)")
	_ = runCMD.MarkFlagRequired("manifest")
	rootCMD.AddCommand(runCMD)
}
