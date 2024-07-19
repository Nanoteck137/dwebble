package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/tools/routes"
	"github.com/nanoteck137/pyrin/ast"
	"github.com/nanoteck137/pyrin/client"
	"github.com/nanoteck137/pyrin/extract"
	"github.com/nanoteck137/pyrin/resolve"
	"github.com/nanoteck137/pyrin/util"
)

// TODO(patrik):
//  - Use json names
//  - Parse omit flag
//  - General code cleanup

type Context struct {
	types    map[string]reflect.Type
	nameUsed map[string]int
	names    map[string]string
}

func (c *Context) RegisterName(name, pkg string) string {
	fullName := pkg + "-" + name

	used, exists := c.nameUsed[name]
	if !exists {
		c.nameUsed[name] = 1
		c.names[fullName] = name

		return name
	} else {
		c.nameUsed[name] = used + 1

		newName := name + strconv.Itoa(used+1)
		c.names[fullName] = newName

		return newName
	}
}

func (c *Context) translateName(name, pkg string) string {
	fullName := pkg + "-" + name

	n, exists := c.names[fullName]
	if !exists {
		log.Fatal("Name not registered", "name", name, "pkg", pkg)
	}

	return n
}

func (c *Context) getType(t reflect.Type) ast.Typespec {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &ast.IdentTypespec{Ident: "int"}
	case reflect.Bool:
		return &ast.IdentTypespec{Ident: "int"}
	case reflect.String:
		return &ast.IdentTypespec{Ident: "string"}
	case reflect.Struct:
		fullName := t.PkgPath() + "-" + t.Name()
		n, exists := c.names[fullName]
		if !exists {
			name := c.RegisterName(t.Name(), t.PkgPath())

			_, exists := c.types[name]
			if !exists {
				c.types[name] = t
			}

			n = name
		}

		return &ast.IdentTypespec{Ident: n}
	case reflect.Slice:
		el := c.getType(t.Elem())
		return &ast.ArrayTypespec{
			Element: el,
		}
	default:
		log.Fatal("Unknown type", "name", t.Name(), "kind", t.Kind())
	}

	return nil
}

func printIndent(indent int) {
	if indent == 0 {
		return
	}

	for i := 0; i < indent; i++ {
		fmt.Print("  ")
	}
}

func checkType(t reflect.Type, indent int) {
	switch t.Kind() {
	case reflect.Struct:
		printStruct(t, indent+1)
	case reflect.Slice:
		checkType(t.Elem(), indent)
	}
}

func printStruct(t reflect.Type, indent int) {
	if t.Kind() != reflect.Struct {
		log.Fatal("Type needs to be struct")
	}

	printIndent(indent)
	fmt.Println("Name: ", t.Name())

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)

		// printIndent(indent)
		// fmt.Printf("sf.Type.Kind(): %v\n", sf.Type.Kind())

		checkType(sf.Type, indent)
	}
}

func main() {
	routes := routes.ServerRoutes(nil)

	s := client.Server{}

	resolver := resolve.New()

	c2 := extract.NewContext()

	for _, route := range routes {
		if route.Data != nil {
			c2.ExtractTypes(route.Data)
		}

		if route.Body != nil {
			c2.ExtractTypes(route.Body)
		}
	}

	pretty.Println(c2)

	decls, err := c2.ConvertToDecls()
	if err != nil {
		log.Fatal("Failed to convert extract context to decls", "err", err)
	}

	pretty.Println(decls)

	for _, decl := range decls {
		resolver.AddSymbolDecl(decl)
	}

	for _, route := range routes {
		responseType := ""
		bodyType := ""

		if route.Data != nil {
			t := reflect.TypeOf(route.Data)

			name, err := c2.TranslateName(t.Name(), t.PkgPath())
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

			name, err := c2.TranslateName(t.Name(), t.PkgPath())
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

	log.Info("Wrote 'pyrin.json'")
}
