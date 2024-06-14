package gen

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/nanoteck137/dwebble/filter"
)

func generateTableSelect(table *filter.Table) (*goqu.SelectDataset, error) {
	switch table.Type {
	case "tags":
		return goqu.From("tracks_to_tags").
			Select("track_id").
			Where(goqu.I("tag_id").In(table.Ids)), nil
	case "genres":
		return goqu.From("tracks_to_genres").
			Select("track_id").
			Where(goqu.I("genre_id").In(table.Ids)), nil
	}

	return nil, fmt.Errorf("Unsupported type: %s\n", table.Type)
}

func Generate(e filter.FilterExpr) (exp.Expression, error) {
	switch e := e.(type) {
	case *filter.AndExpr:
		left, err := Generate(e.Left)
		if err != nil {
			return nil, err
		}

		right, err := Generate(e.Right)
		if err != nil {
			return nil, err
		}

		return goqu.L("(? AND ?)", left, right), nil
	case *filter.OrExpr:
		left, err := Generate(e.Left)
		if err != nil {
			return nil, err
		}

		right, err := Generate(e.Right)
		if err != nil {
			return nil, err
		}

		return goqu.L("(? OR ?)", left, right), nil
	case *filter.OpExpr:
		switch e.Kind {
		case filter.OpEqual:
			return goqu.L("(? == ?)", goqu.I(e.Name), e.Value), nil
		case filter.OpNotEqual:
			return goqu.L("(? != ?)", goqu.I(e.Name), e.Value), nil
		case filter.OpLike:
			return goqu.L("(? LIKE ?)", goqu.I(e.Name), e.Value), nil
		case filter.OpGreater:
			return goqu.L("(? > ?)", goqu.I(e.Name), e.Value), nil
		default:
			return nil, fmt.Errorf("Unimplemented OpKind %d", e.Kind)
		}
	case *filter.InTableExpr:
		s, err := generateTableSelect(&e.Table)
		if err != nil {
			return nil, err
		}

		if e.Not {
			return goqu.L("? NOT IN ?", goqu.I("tracks.id"), s), nil
		} else {
			return goqu.L("? IN ?", goqu.I("tracks.id"), s), nil
		}
	}

	return nil, fmt.Errorf("Unimplemented expr %T", e)
}
