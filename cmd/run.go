package cmd

import (
	"fmt"
	"github.com/mkawserm/abesh/model"
	"github.com/spf13/cobra"
	"os"
)

var runCMD = &cobra.Command{
	Use:   "run",
	Short: "run the platform",
	Long:  "run all platform triggers and services",
	Run: func(cmd *cobra.Command, args []string) {
		if len(manifestFilePath) == 0 {
			fmt.Println("manifest file path required")
			os.Exit(1)
		}

		manifest, err := model.GetManifestFromFile(manifestFilePath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		err = DefaultPlatform.Setup(manifest)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		DefaultPlatform.Run()
	},
}

func init() {
	runCMD.Flags().StringVar(&manifestFilePath, "manifest", "", "Manifest file path (ex: /home/ubuntu/manifest.yaml)")
	_ = runCMD.MarkFlagRequired("manifest")
	rootCMD.AddCommand(runCMD)
}
