package main

import "github.com/mkawserm/abesh/cmd"
import _ "github.com/mkawserm/abesh/capability/httpserver"
import _ "github.com/mkawserm/abesh/capability/httpclient"
import _ "github.com/mkawserm/abesh/example/echo"
import _ "github.com/mkawserm/abesh/example/authorizer"
import _ "github.com/mkawserm/abesh/example/consumer"
import _ "github.com/mkawserm/abesh/example/exhttpclient"

func main() {
	cmd.Execute()
}
