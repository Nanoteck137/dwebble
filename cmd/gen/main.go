package main

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/log"
	"github.com/nanoteck137/dwebble/server"
	"github.com/nanoteck137/pyrin/ast"
	"github.com/nanoteck137/pyrin/client"
	"github.com/nanoteck137/pyrin/parser"
	"github.com/nanoteck137/pyrin/resolve"
	"github.com/nanoteck137/pyrin/util"
)

func main() {
	routes := server.ServerRoutes()

	pretty.Println(routes)

	d, err := os.ReadFile("./types/api_types.go")
	if err != nil {
		log.Fatal("Failed to read api types source", "err", err)
	}

	decls := parser.Parse(string(d))

	findDecl := func(name string) ast.Decl {
		for _, decl := range decls {
			switch decl := decl.(type) {
			case *ast.StructDecl:
				if decl.Name == name {
					return decl
				}
			}
		}

		return nil
	}

	_ = findDecl

	pretty.Println(decls)

	s := client.Server{}

	resolver := resolve.New()

	for _, decl := range decls {
		resolver.AddSymbolDecl(decl)
	}

	for _, route := range routes {
		responseType := ""
		bodyType := ""

		if route.Data != nil {
			t := reflect.TypeOf(route.Data)
			_, err := resolver.Resolve(t.Name())
			if err != nil {
				log.Fatal("Failed to resolve", "name", t.Name(), "err", err)
			}

			responseType = t.Name()
		}

		if route.Body != nil {
			t := reflect.TypeOf(route.Body)
			_, err := resolver.Resolve(t.Name())
			if err != nil {
				log.Fatal("Failed to resolve", "name", t.Name(), "err", err)
			}

			bodyType = t.Name()
		}

		s.Endpoints = append(s.Endpoints, client.Endpoint{
			Name:         route.Name,
			Method:       route.Method,
			Path:         route.Path,
			ResponseType: responseType,
			BodyType:     bodyType,
		})
	}

	pretty.Println(resolver.ResolvedStructs)

	for _, st := range resolver.ResolvedStructs {
		switch t := st.Type.(type) {
		case *resolve.TypeStruct:
			fields := make([]client.TypeField, 0, len(t.Fields))

			for _, f := range t.Fields {
				s, err := util.TypeToString(f.Type)
				if err != nil {
					log.Fatal("TypeToString failed", "err", err)
				}

				fields = append(fields, client.TypeField{
					Name: f.Name,
					Type: s,
					Omit: f.Optional,
				})
			}

			s.Types = append(s.Types, client.Type{
				Name:   st.Name,
				Extend: "",
				Fields: fields,
			})
		case *resolve.TypeSameStruct:
			s.Types = append(s.Types, client.Type{
				Name:   st.Name,
				Extend: t.Type.Name,
			})
		}
	}

	pretty.Println(s)

	d, err = json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal server", "err", err)
	}

	out := "./misc/pyrin.json"
	err = os.WriteFile(out, d, 0644)
	if err != nil {
		log.Fatal("Failed to write pyrin.json", "err", err)
	}
}
