package main

import (
	"embed"
	"github.com/mkawserm/abesh/capability/httpserver"
	"github.com/mkawserm/abesh/cmd"
	"github.com/spf13/cobra"
)
import _ "github.com/mkawserm/abesh/capability/httpserver"
import _ "github.com/mkawserm/abesh/capability/httpclient"
import _ "github.com/mkawserm/abesh/capability/pprof"
import _ "github.com/mkawserm/abesh/capability/health"
import _ "github.com/mkawserm/abesh/example/echo"
import _ "github.com/mkawserm/abesh/example/authorizer"
import _ "github.com/mkawserm/abesh/example/consumer"
import _ "github.com/mkawserm/abesh/example/exhttpclient"
import _ "github.com/mkawserm/abesh/example/exerr"
import _ "github.com/mkawserm/abesh/example/exrpc"
import _ "github.com/mkawserm/abesh/example/expanic"
import _ "embed"

//go:embed manifest.yaml
var manifestBytes []byte

//go:embed data
var staticDataFiles embed.FS

var manifestFilePathList []string

var embeddedRunCMD2 = &cobra.Command{
	Use:   "embedded-run2",
	Short: "Run the platform in embedded mode 2",
	Long:  "Run all platform components with the embedded manifest as source manifest",
	Run: func(c *cobra.Command, args []string) {
		p := cmd.EmbeddedPlatformSetup(manifestFilePathList)
		t := p.GetTriggersCapability()["abesh:httpserver"]
		srv := t.(*httpserver.HTTPServer)
		srv.AddEmbeddedStaticFS("/data/", staticDataFiles)
		p.Run()
	},
}

func main() {
	embeddedRunCMD2.PersistentFlags().StringSliceVar(&manifestFilePathList, "manifest", []string{}, "Manifest file path list (ex: /home/manifest1.yaml,/home/manifest2.yaml)")
	cmd.AddCommand(embeddedRunCMD2)

	cmd.ManifestBytes = manifestBytes
	cmd.Execute()
}
