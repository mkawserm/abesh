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

func EmbeddedPlatformSetup(manifestFilePath string) iface.IPlatform {
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

	if toManifest == nil && fromManifest == nil {
		fmt.Println("no manifest found")
		os.Exit(1)
	}

	if toManifest == nil {
		toManifest = fromManifest
		fromManifest = nil
	}

	currentManifest := utility.MergeManifest(toManifest, fromManifest)
	p := PlatformSetupWithManifest(currentManifest)
	return p
}
