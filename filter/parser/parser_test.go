package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/filter/parser"
)

func TestParser(t *testing.T) {
	src := `has(tag."Anime", genre."Nerdcore") && (not(genre."Soundtrack") || has(tag."One Piece"))`
	p := parser.New(strings.NewReader(src))

	e := p.Parse()
	pretty.Println(e)

	for _, e := range p.Errors {
		fmt.Printf("e: %v\n", e)
	}
}
