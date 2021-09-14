package cmd

import (
	"fmt"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/model"
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

	err = DefaultPlatform.Setup(manifest)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return DefaultPlatform
}
