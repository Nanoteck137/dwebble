package filter

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/nanoteck137/dwebble/types"
)

type IdMappingFunc func(typ string, name string) string

type ResolverAdapter interface {
	MapNameToId(typ, name string) (string, error)
	// TODO(patrik): Rename to ResolveVariableName
	MapName(name string) (types.Name, error)

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

func (r *Resolver) ResolveToIdent(e ast.Expr) (string, error) {
	ident, ok := e.(*ast.Ident)
	if !ok {
		return "", fmt.Errorf("Expected identifier got %T", e)
	}

	return ident.Name, nil
}

func (r *Resolver) ResolveToStr(e ast.Expr) (string, error) {
	lit, ok := e.(*ast.BasicLit)
	if !ok {
		return "", fmt.Errorf("Expected string got %T", e)
	}

	if lit.Kind != token.STRING {
		return "", fmt.Errorf("Expected string got %s", lit.Value)
	}

	s, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", fmt.Errorf("Failed to unquote string: %s", lit.Value)
	}

	return s, nil
}

func (r *Resolver) ResolveToNumber(e ast.Expr) (int64, error) {
	lit, ok := e.(*ast.BasicLit)
	if !ok {
		return 0, fmt.Errorf("Expected number got %T", e)
	}

	if lit.Kind != token.INT {
		return 0, fmt.Errorf("Expected number got %s", lit.Value)
	}

	i, err := strconv.ParseInt(lit.Value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse number: %s", lit.Value)
	}

	return i, nil
}

func (r *Resolver) resolveNameValue(name string, value ast.Expr) (string, any, error) {
	n, err := r.adapter.MapName(name)
	if err != nil {
		return "", nil, err
	}

	var val any

	switch n.Kind {
	case types.NameKindString:
		val, err = r.ResolveToStr(value)
		if err != nil {
			return "", nil, err
		}
	case types.NameKindNumber:
		val, err = r.ResolveToNumber(value)
		if err != nil {
			return "", nil, err
		}
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
		s, err := r.ResolveToStr(arg)
		if err != nil {
			return nil, err
		}

		id, err := r.adapter.MapNameToId(typ, s)
		if err != nil {
			return nil, err
		}

		if id != "" {
			ids = append(ids, id)
		}
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
			name, err := r.ResolveToIdent(e.X)
			if err != nil {
				return nil, err
			}

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
			name, err := r.ResolveToIdent(e.X)
			if err != nil {
				return nil, err
			}

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
			name, err := r.ResolveToIdent(e.X)
			if err != nil {
				return nil, err
			}

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
			name, err := r.ResolveToIdent(e.X)
			if err != nil {
				return nil, err
			}

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
		name, err := r.ResolveToIdent(e.Fun)
		if err != nil {
			return nil, err
		}

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

	switch e := e.(type) {
	case *ast.Ident:
		return nil, fmt.Errorf("Unexpected identifier: %s", e.String())
	case *ast.BasicLit:
		return nil, fmt.Errorf("Unexpected literal: %s", e.Value)
	default:
		return nil, fmt.Errorf("Unexpected expr: %T", e)
	}
}
