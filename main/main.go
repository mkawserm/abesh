package main

import "github.com/mkawserm/abesh/cmd"
import _ "github.com/mkawserm/abesh/capability/httpserver"
import _ "github.com/mkawserm/abesh/example/echo"
import _ "github.com/mkawserm/abesh/example/authorizer"
import _ "github.com/mkawserm/abesh/example/consumer"

func main() {
	cmd.Execute()
}
