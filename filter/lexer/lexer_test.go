package lexer_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nanoteck137/dwebble/filter/lexer"
	"github.com/nanoteck137/dwebble/filter/token"
)

func TestLexer(t *testing.T) {
	src := `"test" hello world ()[]`

	lexer := lexer.New(strings.NewReader(src))

	for {
		tok := lexer.NextToken()
		if tok.Kind == token.Eof {
			break
		}

		fmt.Printf("tok.Kind: %v\n", tok.Kind)
	}
}
