package lexer

import (
	"bufio"
	"bytes"
	"io"
	"unicode"

	"github.com/nanoteck137/dwebble/filter/token"
)

var eof = rune(0)

type Tokenizer struct {
	reader *bufio.Reader
	pos    token.Pos
}

func New(reader io.Reader) *Tokenizer {
	return &Tokenizer{
		reader: bufio.NewReader(reader),
		pos: token.Pos{
			Line:   1,
			Column: 1,
		},
	}
}

func (t *Tokenizer) read() rune {
	ch, _, err := t.reader.ReadRune()
	if err != nil {
		return eof
	}

	if ch == '\n' {
		t.pos.Line += 1
	}

	return ch
}

func (t *Tokenizer) unread() {
	t.reader.UnreadRune()
}

func (t *Tokenizer) NextToken() token.Token {
	c := t.read()

	for unicode.IsSpace(c) {
		c = t.read()
	}

	pos := t.pos

	if unicode.IsLetter(c) || c == '_' {
		var b bytes.Buffer

		for unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' {
			b.WriteRune(c)
			c = t.read()
		}

		t.unread()

		ident := b.String()
		kind := token.Lookup(ident)

		return token.Token{
			Kind:  kind,
			Ident: ident,
			Pos:   pos,
		}
	}

	if c == '"' {
		c = t.read()

		var b bytes.Buffer
		for c != '"' {
			b.WriteRune(c)
			c = t.read()

			if c == eof {
				break
			}
		}

		s := b.String()
		return token.Token{
			Kind:  token.Str,
			Ident: s,
			Pos:   pos,
		}
	}

	kind := token.Unknown
	switch c {
	case eof:
		kind = token.Eof
	case '{':
		kind = token.LBrace
	case '}':
		kind = token.RBrace
	case '[':
		kind = token.LBracket
	case ']':
		kind = token.RBracket
	case '(':
		kind = token.LParen
	case ')':
		kind = token.RParen
	case ',':
		kind = token.Comma
	case '&':
		c = t.read()
		if c == '&' {
			kind = token.DoubleAnd
		} else {
			t.unread()
			kind = token.And
		}
	case '|':
		if c := t.read(); c == '|' {
			kind = token.DoubleOr
		} else {
			t.unread()
			kind = token.Or
		}
	case '=':
		if c := t.read(); c == '=' {
			kind = token.DoubleEqual
		} else {
			t.unread()
			kind = token.Equal
		}
	case '!':
		if c := t.read(); c == '=' {
			kind = token.NotEqual
		} else {
			t.unread()
			kind = token.Unknown
		}
	case '.':
		kind = token.Dot
	}

	return token.Token{
		Kind:  kind,
		Ident: "",
		Pos:   pos,
	}
}
