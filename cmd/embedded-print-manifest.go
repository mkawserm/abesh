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
		var toManifest *model.Manifest
		var currentManifest *model.Manifest

		if len(ManifestBytes) != 0 {
			manifest, err := model.GetManifestFromBytes(ManifestBytes)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			toManifest = manifest
		}

		currentManifest = toManifest

		for _, mfp := range manifestFilePathList {
			if len(mfp) != 0 {
				manifest, err := model.GetManifestFromFile(mfp)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				if currentManifest == nil && manifest != nil {
					currentManifest = manifest
				} else {
					currentManifest = utility.MergeManifest(currentManifest, manifest)
				}
			}
		}

		if currentManifest == nil {
			fmt.Println("no embedded manifest found")
			os.Exit(1)
		}

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
