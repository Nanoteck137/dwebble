package types

import (
	"errors"
	"fmt"
)

type NameKind int

const (
	NameKindString NameKind = iota
	NameKindNumber
)

type Name struct {
	Kind NameKind
	Name string
}

var ErrUnknownName = errors.New("Unknown name")

func UnknownName(name string) error {
	return fmt.Errorf("%w: %s", ErrUnknownName, name)
}
