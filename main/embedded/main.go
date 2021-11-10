package main

import (
	"github.com/mkawserm/abesh/cmd"
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

func main() {
	cmd.ManifestBytes = manifestBytes
	cmd.Execute()
}
