package filter

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"strconv"
)

type IdMappingFunc func(typ string, name string) string

type ResolverAdapter interface {
	GetDefaultSort() string

	MapNameToId(typ, name string) (string, error)
	// TODO(patrik): Rename to ResolveVariableName
	MapName(name string) (Name, error)

	ResolveTable(typ string) (Table, error)

	ResolveFunctionCall(resolver *Resolver, name string, args []ast.Expr) (FilterExpr, error)
}

type Resolver struct {
	adapter ResolverAdapter
}

func New(adpater ResolverAdapter) *Resolver {
	return &Resolver{
		adapter: adpater,
	}
}

func (r *Resolver) ResolveToIdent(e ast.Expr) string {
	ident, ok := e.(*ast.Ident)
	if !ok {
		// TODO(patrik): Return error
		panic("Expected ident")
	}

	return ident.Name
}

func (r *Resolver) ResolveToStr(e ast.Expr) string {
	lit, ok := e.(*ast.BasicLit)
	if !ok {
		// TODO(patrik): Return error
		panic("Expected BasicLit")
	}

	if lit.Kind != token.STRING {
		// TODO(patrik): Return error
		panic("Expected string")
	}

	s, err := strconv.Unquote(lit.Value)
	if err != nil {
		// TODO(patrik): Return err
		log.Fatal(err)
	}

	return s
}

func (r *Resolver) ResolveToNumber(e ast.Expr) int64 {
	lit, ok := e.(*ast.BasicLit)
	if !ok {
		// TODO(patrik): Return error
		panic("Expected BasicLit")
	}

	if lit.Kind != token.INT {
		// TODO(patrik): Return error
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

func (r *Resolver) resolveNameValue(name string, value ast.Expr) (string, any, error) {
	n, err := r.adapter.MapName(name)
	if err != nil {
		return "", nil, err
	}

	var val any

	switch n.Kind {
	case NameKindString:
		val = r.ResolveToStr(value)
	case NameKindNumber:
		val = r.ResolveToNumber(value)
	default:
		return "", nil, fmt.Errorf("Unknown NameKind: %d", n.Kind)
	}

	return n.Name, val, nil
}

func (r *Resolver) InTable(name, typ string, args []ast.Expr) (*InTableExpr, error) {
	if len(args) <= 0 {
		return nil, fmt.Errorf("'%s' requires at least 1 parameter", name)
	}

	var ids []string
	for _, arg := range args {
		s := r.ResolveToStr(arg)
		id, err := r.adapter.MapNameToId(typ, s)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	tbl, err := r.adapter.ResolveTable(typ)
	if err != nil {
		return nil, err
	}

	return &InTableExpr{
		Not:   false,
		Table: tbl,
		Ids:   ids,
	}, nil
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
			name := r.ResolveToIdent(e.X)
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
			name := r.ResolveToIdent(e.X)
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
			name := r.ResolveToIdent(e.X)
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
			name := r.ResolveToIdent(e.X)
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
		name := r.ResolveToIdent(e.Fun)
		return r.adapter.ResolveFunctionCall(r, name, e.Args)
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

	// TODO(patrik): Return error
	panic(fmt.Sprintf("Unexpected expr: %T", e))
}
