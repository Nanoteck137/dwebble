package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/filter/parser"
)

func TestParser(t *testing.T) {
	// has(tag."Hello" | genres."Soundtrack") && not(tag."Anime")

	// src := `&& tags["Hello"] + genres["Soundtrack"] &~ tags["Anime"]`
	src := `has(tag."Anime") && not(genre."Soundtrack") || has(tag."One Piece")`
	p := parser.New(strings.NewReader(src))

	expr := p.ParseExpr()
	pretty.Println(expr)

	expr = p.ParseExpr()
	pretty.Println(expr)

	for _, v := range p.Errors {
		fmt.Printf("v: %v\n", v)
	}
}
