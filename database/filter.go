package database

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/nanoteck137/dwebble/tools/filter"
	"github.com/nanoteck137/dwebble/tools/sort"
)

func generateTableSelect(table *filter.Table, ids []string) *goqu.SelectDataset {
	return goqu.From(table.Name).
		Select(table.SelectName).
		Where(goqu.I(table.WhereName).In(ids))
}

func generateFilter(e filter.FilterExpr) (exp.Expression, error) {
	switch e := e.(type) {
	case *filter.AndExpr:
		left, err := generateFilter(e.Left)
		if err != nil {
			return nil, err
		}

		right, err := generateFilter(e.Right)
		if err != nil {
			return nil, err
		}

		return goqu.L("(? AND ?)", left, right), nil
	case *filter.OrExpr:
		left, err := generateFilter(e.Left)
		if err != nil {
			return nil, err
		}

		right, err := generateFilter(e.Right)
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
		s := generateTableSelect(&e.Table, e.Ids)

		if e.Not {
			// TODO(patrik): Move tracks.id
			return goqu.L("? NOT IN ?", goqu.I("tracks.id"), s), nil
		} else {
			return goqu.L("? IN ?", goqu.I("tracks.id"), s), nil
		}
	}

	return nil, fmt.Errorf("Unimplemented expr %T", e)
}

func generateSort(e sort.SortExpr) ([]exp.OrderedExpression, error) {
	switch e := e.(type) {
	case *sort.SortExprSort:
		var items []exp.OrderedExpression
		for _, item := range e.Items {
			switch item.Type {
			case sort.SortTypeAsc:
				items = append(items, goqu.I(item.Name).Asc())
			case sort.SortTypeDesc:
				items = append(items, goqu.I(item.Name).Desc())
			default:
				return nil, fmt.Errorf("Unknown SortItemType: %d", item.Type)
			}
		}

		return items, nil
	case *sort.SortExprRandom:
		oe := goqu.Func("RANDOM").Asc()
		return []exp.OrderedExpression{oe}, nil
	}

	return nil, fmt.Errorf("Unimplemented expr %T", e)
}
