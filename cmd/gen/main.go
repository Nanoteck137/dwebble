package main

import (
	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/server"
	"github.com/nanoteck137/pyrin/client"
)

func main() {
	routes := server.ServerRoutes()

	pretty.Println(routes)

	s := client.Server{}

	for _, route := range routes {
		s.Endpoints = append(s.Endpoints, client.Endpoint{
			Name:         route.Name,
			Method:       route.Method,
			Path:         route.Path,
			ResponseType: "",
			BodyType:     "",
		})
	}

	pretty.Println(s)
}
