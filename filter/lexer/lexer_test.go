package lexer_test

import (
	"strings"
	"testing"

	"github.com/nanoteck137/dwebble/filter/lexer"
	"github.com/nanoteck137/dwebble/filter/token"
)

func assert(t *testing.T, cond bool, format string, args ...any) {
	if !cond {
		t.Fatalf(format, args...)
	}
}

type wrapper struct {
	*lexer.Tokenizer

	current token.Token
}

func (w *wrapper) next() {
	w.current = w.NextToken()
}

func (w *wrapper) expect(t *testing.T, tok token.Kind) {
	assert(t, w.current.Kind == tok, "Expected token %s got %s", tok.String(), w.current.Kind.String())
	w.next()
}

func (w *wrapper) expectIdent(t *testing.T, ident string) {
	i := w.current.Ident
	w.expect(t, token.Ident)

	assert(t, i == ident, "Expected ident '%s' got '%s'", ident, i)
}

func TestLexerIdents(t *testing.T) {
	src := "hello hello_world _bye_world test123 test_123"
	w := wrapper{Tokenizer: lexer.New(strings.NewReader(src))}
	w.next()

	w.expectIdent(t, "hello")
	w.expectIdent(t, "hello_world")
	w.expectIdent(t, "_bye_world")
	w.expectIdent(t, "test123")
	w.expectIdent(t, "test_123")
}
