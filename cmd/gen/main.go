package main

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/tools/routes"
	"github.com/nanoteck137/pyrin/client"
	"github.com/nanoteck137/pyrin/extract"
	"github.com/nanoteck137/pyrin/resolve"
	"github.com/nanoteck137/pyrin/util"
)

// TODO(patrik):
//  - Use json names
//  - Parse omit flag
//  - General code cleanup

func main() {
	routes := routes.ServerRoutes(nil)

	s := client.Server{}

	resolver := resolve.New()

	c := extract.NewContext()

	for _, route := range routes {
		if route.Data != nil {
			c.ExtractTypes(route.Data)
		}

		if route.Body != nil {
			c.ExtractTypes(route.Body)
		}
	}

	decls, err := c.ConvertToDecls()
	if err != nil {
		log.Fatal("Failed to convert extract context to decls", "err", err)
	}

	for _, decl := range decls {
		resolver.AddSymbolDecl(decl)
	}

	for _, route := range routes {
		responseType := ""
		bodyType := ""

		if route.Data != nil {
			t := reflect.TypeOf(route.Data)

			name, err := c.TranslateName(t.Name(), t.PkgPath())
			if err != nil {
				log.Fatal("Failed to translate name", "name", t.Name(), "pkg", t.PkgPath(), "err", err)
			}

			_, err = resolver.Resolve(name)
			if err != nil {
				log.Fatal("Failed to resolve", "name", t.Name(), "err", err)
			}

			responseType = name
		}

		if route.Body != nil {
			t := reflect.TypeOf(route.Body)

			name, err := c.TranslateName(t.Name(), t.PkgPath())
			if err != nil {
				log.Fatal("Failed to translate name", "name", t.Name(), "pkg", t.PkgPath(), "err", err)
			}

			_, err = resolver.Resolve(name)
			if err != nil {
				log.Fatal("Failed to resolve", "name", t.Name(), "err", err)
			}

			responseType = name
		}

		s.Endpoints = append(s.Endpoints, client.Endpoint{
			Name:         route.Name,
			Method:       route.Method,
			Path:         route.Path,
			ResponseType: responseType,
			BodyType:     bodyType,
		})
	}

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

	d, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal server", "err", err)
	}

	out := "./misc/pyrin.json"
	err = os.WriteFile(out, d, 0644)
	if err != nil {
		log.Fatal("Failed to write pyrin.json", "err", err)
	}

	log.Info("Wrote 'misc/pyrin.json'")
}
