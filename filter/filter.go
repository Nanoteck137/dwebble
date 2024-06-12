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

func (r *Resolver) Resolve(e ast.Expr) (FilterExpr, error) {
	switch e := e.(type) {
	case *ast.BinaryExpr:
		left, err := r.Resolve(e.X)
		if err != nil {
			return nil, err
		}

		right, err := r.Resolve(e.Y)
		if err != nil {
			return nil, err
		}

		switch e.Op {
		case token.LAND:
			return &AndExpr{
				Left:  left,
				Right: right,
			}, nil
		case token.LOR:
			return &OrExpr{
				Left:  left,
				Right: right,
			}, nil
		default:
			return nil, fmt.Errorf("Unsupported binary operator: %s", e.Op.String())
		}
	case *ast.CallExpr:
		name := r.resolveToIdent(e.Fun)
		fmt.Printf("name: %v\n", name)
		switch name {
		case "hasTag":
			expr := e.Args[0]
			s := r.resolveToStr(expr)

			id := r.mapNameToId("tags", s)
			return &InTableExpr{
				Not: false,
				Table: Table{
					Type: "tags",
					Ids:  []string{id},
				},
			}, nil
		case "hasGenre":
			expr := e.Args[0]
			s := r.resolveToStr(expr)

			id := r.mapNameToId("genres", s)
			return &InTableExpr{
				Not: false,
				Table: Table{
					Type: "genres",
					Ids:  []string{id},
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
