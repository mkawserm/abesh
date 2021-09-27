package cmd

import (
	"fmt"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/utility"
	"os"
)

func PlatformSetup(manifestFilePath string) iface.IPlatform {
	if len(manifestFilePath) == 0 {
		fmt.Println("manifest file path required")
		os.Exit(1)
	}

	manifest, err := model.GetManifestFromFile(manifestFilePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return PlatformSetupWithManifest(manifest)
}

func PlatformSetupWithManifest(manifest *model.Manifest) iface.IPlatform {
	err := DefaultPlatform.Setup(manifest)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return DefaultPlatform
}

func EmbeddedPlatformSetup(manifestFilePathList []string) iface.IPlatform {
	//var fromManifest *model.Manifest
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
		fmt.Println("no manifest found")
		os.Exit(1)
	}

	p := PlatformSetupWithManifest(currentManifest)
	return p
}
