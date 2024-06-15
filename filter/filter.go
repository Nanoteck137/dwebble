package filter

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"strconv"
)

type FilterExpr interface {
	filterExprType()
}

type AndExpr struct {
	Left  FilterExpr
	Right FilterExpr
}

type OrExpr struct {
	Left  FilterExpr
	Right FilterExpr
}

type OpKind int

const (
	OpEqual OpKind = iota
	OpNotEqual
	OpLike
	OpGreater
)

type OpExpr struct {
	Kind  OpKind
	Name  string
	Value any
}

type Table struct {
	Type string
	Ids  []string
}

type InTableExpr struct {
	Not   bool
	Table Table
}

func (e *AndExpr) filterExprType()     {}
func (e *OrExpr) filterExprType()      {}
func (e *OpExpr) filterExprType()      {}
func (e *InTableExpr) filterExprType() {}

type IdMappingFunc func(typ string, name string) string

type Resolver struct {
	mapNameToId IdMappingFunc
}

func New(mapNametoId IdMappingFunc) *Resolver {
	return &Resolver{
		mapNameToId: mapNametoId,
	}
}

func (r *Resolver) resolveToIdent(e ast.Expr) string {
	ident, ok := e.(*ast.Ident)
	if !ok {
		panic("Expected ident")
	}

	return ident.Name
}

func (r *Resolver) resolveToStr(e ast.Expr) string {
	lit, ok := e.(*ast.BasicLit)
	if !ok {
		panic("Expected BasicLit")
	}

	if lit.Kind != token.STRING {
		panic("Expected string")
	}

	s, err := strconv.Unquote(lit.Value)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

func (r *Resolver) resolveToNumber(e ast.Expr) int64 {
	lit, ok := e.(*ast.BasicLit)
	if !ok {
		panic("Expected BasicLit")
	}

	if lit.Kind != token.INT {
		panic("Expected int")
	}

	i, err := strconv.ParseInt(lit.Value, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	return i
}

type NameKind int

const (
	NameKindString NameKind = iota
	NameKindNumber
)

type Name struct {
	Kind NameKind
	Name string
}

var GlobalNames = map[string]Name{
	"artist": {
		Kind: NameKindString,
		Name: "artists.name",
	},
	"album": {
		Kind: NameKindString,
		Name: "albums.name",
	},
	"track": {
		Kind: NameKindString,
		Name: "tracks.name",
	},
	"trackNumber": {
		Kind: NameKindNumber,
		Name: "tracks.track_number",
	},
}

func (r *Resolver) resolveNameValue(name string, value ast.Expr) (string, any, error) {
	n, exists := GlobalNames[name]
	if !exists {
		return "", nil, fmt.Errorf("Unknown name: %s", name)
	}

	var val any

	switch n.Kind {
	case NameKindString:
		val = r.resolveToStr(value)
	case NameKindNumber:
		val = r.resolveToNumber(value)
	default:
		panic("Unimplemented NameKind")
	}

	return n.Name, val, nil
}

func (r *Resolver) Resolve(e ast.Expr) (FilterExpr, error) {
	switch e := e.(type) {
	case *ast.BinaryExpr:
		switch e.Op {
		case token.LAND:
			left, err := r.Resolve(e.X)
			if err != nil {
				return nil, err
			}

			right, err := r.Resolve(e.Y)
			if err != nil {
				return nil, err
			}

			return &AndExpr{
				Left:  left,
				Right: right,
			}, nil
		case token.LOR:
			left, err := r.Resolve(e.X)
			if err != nil {
				return nil, err
			}

			right, err := r.Resolve(e.Y)
			if err != nil {
				return nil, err
			}

			return &OrExpr{
				Left:  left,
				Right: right,
			}, nil
		case token.EQL:
			name := r.resolveToIdent(e.X)
			name, value, err := r.resolveNameValue(name, e.Y)
			if err != nil {
				return nil, err
			}

			return &OpExpr{
				Kind:  OpEqual,
				Name:  name,
				Value: value,
			}, nil
		case token.NEQ:
			name := r.resolveToIdent(e.X)
			name, value, err := r.resolveNameValue(name, e.Y)
			if err != nil {
				return nil, err
			}

			return &OpExpr{
				Kind:  OpNotEqual,
				Name:  name,
				Value: value,
			}, nil
		case token.REM:
			name := r.resolveToIdent(e.X)
			name, value, err := r.resolveNameValue(name, e.Y)
			if err != nil {
				return nil, err
			}

			return &OpExpr{
				Kind:  OpLike,
				Name:  name,
				Value: value,
			}, nil
		case token.GTR:
			name := r.resolveToIdent(e.X)
			name, value, err := r.resolveNameValue(name, e.Y)
			if err != nil {
				return nil, err
			}

			return &OpExpr{
				Kind:  OpGreater,
				Name:  name,
				Value: value,
			}, nil
		default:
			return nil, fmt.Errorf("Unsupported binary operator: %s", e.Op.String())
		}
	case *ast.CallExpr:
		name := r.resolveToIdent(e.Fun)
		fmt.Printf("name: %v\n", name)
		switch name {
		case "hasTag":
			if len(e.Args) <= 0 {
				return nil, fmt.Errorf("'hasTag' requires at least 1 parameter")
			}

			var ids []string
			for _, arg := range e.Args {
				s := r.resolveToStr(arg)
				id := r.mapNameToId("tags", s)
				ids = append(ids, id)
			}

			return &InTableExpr{
				Not: false,
				Table: Table{
					Type: "tags",
					Ids:  ids,
				},
			}, nil
		case "hasGenre":
			if len(e.Args) <= 0 {
				return nil, fmt.Errorf("'hasGenre' requires at least 1 parameter")
			}

			var ids []string
			for _, arg := range e.Args {
				s := r.resolveToStr(arg)
				id := r.mapNameToId("genres", s)
				ids = append(ids, id)
			}

			return &InTableExpr{
				Not: false,
				Table: Table{
					Type: "genres",
					Ids:  ids,
				},
			}, nil
		default:
			return nil, fmt.Errorf("Unknown function name: %s", name)
		}
	case *ast.UnaryExpr:
		expr, err := r.Resolve(e.X)
		if err != nil {
			return nil, err
		}

		switch expr := expr.(type) {
		case *InTableExpr:
			expr.Not = true
		}

		return expr, nil
	}

	panic(fmt.Sprintf("Unexpected expr: %T", e))
}
