package types

type NameKind int

const (
	NameKindString NameKind = iota
	NameKindNumber
)

type Name struct {
	Kind NameKind
	Name string
}
