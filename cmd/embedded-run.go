package cmd

import (
	"github.com/spf13/cobra"
)

var embeddedRunCMD = &cobra.Command{
	Use:   "run",
	Short: "Run the platform in embedded mode",
	Long:  "Run all platform components with the embedded manifest as source manifest",
	Run: func(cmd *cobra.Command, args []string) {
		p := EmbeddedPlatformSetup(manifestFilePath)
		p.Run()
	},
}

func init() {
	embeddedCMD.AddCommand(embeddedRunCMD)
}
