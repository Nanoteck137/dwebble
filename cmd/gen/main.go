package main

import (
	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/server"
)

func main() {
	routes := server.ServerRoutes()

	pretty.Println(routes)
}
