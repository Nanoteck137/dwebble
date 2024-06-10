package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/filter/parser"
)

func TestParser(t *testing.T) {
	src := `artist == "test" && has(tag."Anime", genre."Nerdcore") && (not(genre."Soundtrack") || has(tag."One Piece"))`
	p := parser.New(strings.NewReader(src))

	expr := p.ParseExpr()
	pretty.Println(expr)

	expr = p.ParseExpr()
	pretty.Println(expr)

	for _, v := range p.Errors {
		fmt.Printf("v: %v\n", v)
	}
}
