package resolve

import (
	"github.com/nanoteck137/dwebble/filter/ast"
	"github.com/nanoteck137/dwebble/types"
)

type IdMappingFunc func(typ string, name string) string

type Resolver struct {
	mapNameToId IdMappingFunc
	Errors types.ErrorList
}

func New(mapNameToId IdMappingFunc) *Resolver {
	return &Resolver{
		mapNameToId: mapNameToId,
		Errors:      []*types.Error{},
	}
}

func (r *Resolver) error(message string) {
	r.Errors.Add(message)
}

func (r *Resolver) resolveInTable(not bool, params []ast.Expr) *ast.InTableExpr {
	if len(params) != 1 {
		panic("One param")
	}

	p := params[0]

	e, ok := p.(*ast.FieldExpr)
	if !ok {
		panic("Params should only be 'ast.FieldExpr'")
	}

	if e.Ident != "tags" && e.Ident != "genres" {
		// TODO(patrik): Better error message
		r.error("Unexpected field ident: " + e.Ident)
		return nil
	}

	id := r.mapNameToId(e.Ident, e.Field)

	tbl := ast.Table{
		Type: e.Ident,
		Ids:  []string{id},
	}

	return &ast.InTableExpr{
		Not:   not,
		Table: tbl,
	}
}

func (r *Resolver) resolveExpr(e ast.Expr) ast.Expr {
	switch e := e.(type) {
	case *ast.AndExpr:
		e.Left = r.resolveExpr(e.Left)
		e.Right = r.resolveExpr(e.Right)
	case *ast.OrExpr:
		e.Left = r.resolveExpr(e.Left)
		e.Right = r.resolveExpr(e.Right)
	case *ast.CallExpr:
		switch e.Name {
		case "has":
			return r.resolveInTable(false, e.Params)
		case "not":
			return r.resolveInTable(true, e.Params)
		}
	}

	return e
}

func (r *Resolver) Resolve(e ast.Expr) ast.Expr {
	return r.resolveExpr(e)
}
