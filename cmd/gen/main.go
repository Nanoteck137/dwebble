package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/server"
	"github.com/nanoteck137/pyrin/ast"
	"github.com/nanoteck137/pyrin/client"
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

func main() {
	routes := server.ServerRoutes(nil)

	// d, err := os.ReadFile("./types/api_types.go")
	// if err != nil {
	// 	log.Fatal("Failed to read api types source", "err", err)
	// }
	//
	// decls := parser.Parse(string(d))

	s := client.Server{}

	resolver := resolve.New()

	// for _, decl := range decls {
	// 	resolver.AddSymbolDecl(decl)
	// }

	c := Context{
		types:    map[string]reflect.Type{},
		nameUsed: map[string]int{},
		names:    map[string]string{},
	}

	_ = c

	// for _, route := range routes {
	// 	if route.Data != nil {
	// 		t := reflect.TypeOf(route.Data)
	// 		pretty.Println(t.String())
	//
	// 		fullName := t.PkgPath() + "-" + t.Name()
	// 		fmt.Printf("fullName: %v\n", fullName)
	//
	// 		name := t.Name()
	//
	// 		used, exists := c.nameUsed[name]
	// 		if !exists {
	// 			c.nameUsed[name] = 1
	// 			c.names[fullName] = name
	// 		} else {
	// 			c.nameUsed[name] = used + 1
	// 			c.names[fullName] = name + strconv.Itoa(used+1)
	// 		}
	//
	// 		// name, exists := c.names[fullName]
	// 	}
	// }

	for _, route := range routes {
		if route.Data != nil {
			t := reflect.TypeOf(route.Data)
			pretty.Println(t.String())

			fmt.Printf("t.PkgPath(): %v\n", t.PkgPath())

			if t.Kind() != reflect.Struct {
				log.Fatal("Route data need to be struct", "name", route.Name)
			}

			s := ast.StructDecl{}
			s.Name = c.RegisterName(t.Name(), t.PkgPath())

			for i := 0; i < t.NumField(); i++ {
				sf := t.Field(i)

				if sf.Type.Kind() == reflect.Struct {
					s.Extend = sf.Name
					continue
				}

				s.Fields = append(s.Fields, &ast.Field{
					Name: sf.Name,
					Type: c.getType(sf.Type),
					Omit: false,
				})
			}

			resolver.AddSymbolDecl(&s)
		}

		if route.Body != nil {
			t := reflect.TypeOf(route.Body)
			pretty.Println(t.String())

			fmt.Printf("t.PkgPath(): %v\n", t.PkgPath())

			if t.Kind() != reflect.Struct {
				log.Fatal("Route data need to be struct", "name", route.Name)
			}

			s := ast.StructDecl{}
			s.Name = c.RegisterName(t.Name(), t.PkgPath())

			for i := 0; i < t.NumField(); i++ {
				sf := t.Field(i)

				if sf.Type.Kind() == reflect.Struct {
					s.Extend = sf.Name
					continue
				}

				s.Fields = append(s.Fields, &ast.Field{
					Name: sf.Name,
					Type: c.getType(sf.Type),
					Omit: false,
				})
			}

			resolver.AddSymbolDecl(&s)
		}
	}

	for _, t := range c.types {
		s := ast.StructDecl{}

		s.Name = c.translateName(t.Name(), t.PkgPath())

		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)

			if sf.Type.Kind() == reflect.Struct {
				s.Extend = sf.Name
				continue
			}

			s.Fields = append(s.Fields, &ast.Field{
				Name: sf.Name,
				Type: c.getType(sf.Type),
				Omit: false,
			})
		}

		resolver.AddSymbolDecl(&s)
	}

	pretty.Println(c)
	//
	// pretty.Println(resolver)

	for _, route := range routes {
		responseType := ""
		bodyType := ""

		if route.Data != nil {
			t := reflect.TypeOf(route.Data)

			name := c.translateName(t.Name(), t.PkgPath())

			_, err := resolver.Resolve(name)
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
