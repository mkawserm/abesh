package cmd

import (
	"fmt"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/utility"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
)

var embeddedPrintManifestCMD = &cobra.Command{
	Use:   "print-manifest",
	Short: "Print manifest",
	Long:  "Print currently used manifest combined with the provided manifest",
	Run: func(cmd *cobra.Command, args []string) {
		var fromManifest *model.Manifest
		var toManifest *model.Manifest

		if len(manifestFilePath) != 0 {
			manifest, err := model.GetManifestFromFile(manifestFilePath)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fromManifest = manifest
		}

		if len(ManifestBytes) != 0 {
			manifest, err := model.GetManifestFromBytes(ManifestBytes)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			toManifest = manifest
		}

		if toManifest == nil {
			fmt.Println("no embedded manifest found")
			os.Exit(1)
		}

		currentManifest := utility.MergeManifest(toManifest, fromManifest)

		d, err := yaml.Marshal(currentManifest)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("----------------------")
		fmt.Println(string(d))
		fmt.Println("----------------------")
	},
}

func init() {
	embeddedCMD.AddCommand(embeddedPrintManifestCMD)
}
