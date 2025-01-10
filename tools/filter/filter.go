package filter

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

type NameKind int

const (
	NameKindString NameKind = iota
	NameKindNumber
)

type Name struct {
	Kind     NameKind
	Name     string
	Nullable bool
}

var ErrInternalError = errors.New("internal error")
var ErrUnknownName = errors.New("unknown name")
var ErrUnknownFunction = errors.New("unknown function")

func UnknownName(name string) error {
	return fmt.Errorf("%w: %s", ErrUnknownName, name)
}

func UnknownFunction(name string) error {
	return fmt.Errorf("%w: %s", ErrUnknownFunction, name)
}

func InternalError(err error) error {
	return fmt.Errorf("%w: %w", ErrInternalError, err)
}

type ResolverAdapter interface {
	ResolveNameToId(typ, name string) (string, bool)
	ResolveVariableName(name string) (Name, bool)
	ResolveTable(typ string) (Table, bool)

	ResolveFunctionCall(resolver *Resolver, name string, args []ast.Expr) (FilterExpr, error)

	DefaultSort() (string, SortType)
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
		switch e := e.(type) {
		case *ast.BasicLit:
			return "", fmt.Errorf("expected identifier got %v", e.Value)
		default:
			return "", InternalError(fmt.Errorf("expected identifier got %T", e))
		}

	}

	return ident.Name, nil
}

func (r *Resolver) ResolveToStr(e ast.Expr) (string, error) {
	lit, ok := e.(*ast.BasicLit)
	if !ok {
		switch e := e.(type) {
		case *ast.Ident:
			return "", fmt.Errorf("expected string got %v", e.Name)
		default:
			return "", InternalError(fmt.Errorf("expected string got %T", e))
		}
	}

	if lit.Kind != token.STRING {
		return "", fmt.Errorf("expected string got %s", lit.Value)
	}

	s, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", fmt.Errorf("failed to unquote string %s", lit.Value)
	}

	return s, nil
}

func (r *Resolver) ResolveToNumber(e ast.Expr) (int64, error) {
	lit, ok := e.(*ast.BasicLit)
	if !ok {
		switch e := e.(type) {
		case *ast.Ident:
			return 0, fmt.Errorf("expected number got %v", e.Name)
		default:
			return 0, InternalError(fmt.Errorf("expected number got %T", e))
		}
	}

	if lit.Kind != token.INT {
		return 0, fmt.Errorf("expected number got %s", lit.Value)
	}

	i, err := strconv.ParseInt(lit.Value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse number %s", lit.Value)
	}

	return i, nil
}

func (r *Resolver) resolveNameValue(name string, value ast.Expr) (string, any, error) {
	n, ok := r.adapter.ResolveVariableName(name)
	if !ok {
		return "", nil, UnknownName(name)
	}

	var err error
	var val any

	switch n.Kind {
	case NameKindString:
		val, err = r.ResolveToStr(value)
		if err != nil {
			return "", nil, err
		}
	case NameKindNumber:
		val, err = r.ResolveToNumber(value)
		if err != nil {
			return "", nil, err
		}
	default:
		return "", nil, InternalError(fmt.Errorf("unknown name kind %d", n.Kind))
	}

	return n.Name, val, nil
}

func (r *Resolver) InTable(name, typ, idSelector string, args []ast.Expr) (*InTableExpr, error) {
	if len(args) <= 0 {
		return nil, fmt.Errorf("'%s' requires at least 1 parameter", name)
	}

	var ids []string
	for _, arg := range args {
		s, err := r.ResolveToStr(arg)
		if err != nil {
			return nil, err
		}

		// TODO(patrik): Look at the error here
		id, ok := r.adapter.ResolveNameToId(typ, s)
		if !ok {
			return nil, UnknownName(s)
		}

		if id != "" {
			ids = append(ids, id)
		}
	}

	tbl, ok := r.adapter.ResolveTable(typ)
	if !ok {
		// TODO(patrik): Create custom error here, this might also
		// be an internal error
		return nil, UnknownName(typ)
	}

	return &InTableExpr{
		Not:        false,
		IdSelector: idSelector,
		Table:      tbl,
		Ids:        ids,
	}, nil
}

var opMapping = map[token.Token]OpKind{
	token.EQL: OpEqual,
	token.NEQ: OpNotEqual,
	token.REM: OpLike,
	token.GTR: OpGreater,
	token.GEQ: OpGreaterEqual,
	token.LSS: OpLesser,
	token.LEQ: OpLesserEqual,
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
		default:
			if op, ok := opMapping[e.Op]; ok {
				name, err := r.ResolveToIdent(e.X)
				if err != nil {
					return nil, err
				}

				if ident, ok := e.Y.(*ast.Ident); ok {
					if ident.Name == "null" {
						n, ok := r.adapter.ResolveVariableName(name)
						if !ok {
							return nil, UnknownName(name)
						}

						if !n.Nullable {
							return nil, fmt.Errorf("%s is not nullable", name)
						}

						return &IsNullExpr{
							Name: n.Name,
							Not:  e.Op != token.EQL,
						}, nil
					}
				}

				name, value, err := r.resolveNameValue(name, e.Y)
				if err != nil {
					return nil, err
				}

				return &OpExpr{
					Kind:  op,
					Name:  name,
					Value: value,
				}, nil
			}

			return nil, InternalError(fmt.Errorf("unsupported binary operator %s", e.Op.String()))
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
		return nil, fmt.Errorf("unexpected identifier %s", e.String())
	case *ast.BasicLit:
		return nil, fmt.Errorf("unexpected literal %s", e.Value)
	default:
		return nil, InternalError(fmt.Errorf("unexpected expr %T", e))
	}
}
