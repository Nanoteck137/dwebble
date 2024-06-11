package gen

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/nanoteck137/dwebble/filter/ast"
)

func resolveTable(table *ast.Table) (*goqu.SelectDataset, error) {
	switch table.Type {
	case "tags":
		return goqu.From("tracks_to_tags").Select("track_id").Where(goqu.I("tag_id").In(table.Ids)), nil
	case "genres":
		return goqu.From("tracks_to_genres").Select("track_id").Where(goqu.I("genre_id").In(table.Ids)), nil
	}

	return nil, fmt.Errorf("Unsupported type: %s\n", table.Type)
}

func Generate(e ast.Expr) (exp.Expression, error) {
	switch e := e.(type) {
	case *ast.AndExpr:
		left, err := Generate(e.Left)
		if err != nil {
			return nil, err
		}

		right, err := Generate(e.Right)
		if err != nil {
			return nil, err
		}

		return goqu.L("(? AND ?)", left, right), nil
	case *ast.OrExpr:
		left, err := Generate(e.Left)
		if err != nil {
			return nil, err
		}

		right, err := Generate(e.Right)
		if err != nil {
			return nil, err
		}

		return goqu.L("(? OR ?)", left, right), nil
	case *ast.InTableExpr:
		s, err := resolveTable(&e.Table)
		if err != nil {
			return nil, err
		}

		if e.Not {
			return goqu.L("? NOT IN ?", goqu.I("tracks.id"), s), nil
		} else {
			return goqu.L("? IN ?", goqu.I("tracks.id"), s), nil
		}
	}

	return nil, errors.New("Unexpected expr: " + reflect.TypeOf(e).Name())
}
