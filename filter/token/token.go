package token

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

	Semicolon
	Colon
	DoubleColon
	AndNot
	And
	DoubleAnd
	Or
	DoubleOr

	Asterisk
	Dot
	Plus

	// Unset

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

func Lookup(ident string) Kind {
	// if ident == "unset" {
	// 	return Unset
	// }

	return Ident
}
