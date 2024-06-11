package token

import "strconv"

type Kind int

const (
	Unknown Kind = iota

	Ident
	Str

	LBrace
	RBrace
	LBracket
	RBracket
	LParen
	RParen

	And
	DoubleAnd
	Or
	DoubleOr
	Equal
	DoubleEqual
	NotEqual

	Comma
	Dot

	Eof
)

type Pos struct {
	Line   int
	Column int
}

type Token struct {
	Kind  Kind
	Ident string
	Pos   Pos
}

var tokens = [...]string{
	Unknown: "Unknown",

	Ident: "Ident",
	Str:   "Str",

	LBrace:   "{",
	RBrace:   "}",
	LBracket: "[",
	RBracket: "]",
	LParen:   "(",
	RParen:   ")",

	And:         "&",
	DoubleAnd:   "&&",
	Or:          "|",
	DoubleOr:    "||",
	Equal:       "=",
	DoubleEqual: "==",
	NotEqual:    "!=",

	Comma: ",",
	Dot:   ".",

	Eof: "Eof",
}

func (tok Kind) String() string {
	if 0 <= tok && tok < Kind(len(tokens)) {
		return tokens[tok]
	}

	return "token(" + strconv.Itoa(int(tok)) + ")"
}

func Lookup(ident string) Kind {
	// if ident == "unset" {
	// 	return Unset
	// }

	return Ident
}
