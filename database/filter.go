package database

import (
	"errors"
	"fmt"
	"go/parser"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/nanoteck137/dwebble/tools/filter"
)

// TODO(patrik): Move
var ErrInvalidFilter = errors.New("invalid filter")
var ErrInvalidSort = errors.New("invalid sort")

func InvalidFilter(err error) error {
	return fmt.Errorf("%w: %w", ErrInvalidFilter, err)
}

func InvalidSort(err error) error {
	return fmt.Errorf("%w: %w", ErrInvalidSort, err)
}

func applyFilter(query *goqu.SelectDataset, resolver *filter.Resolver, filterStr string) (*goqu.SelectDataset, error) {
	if filterStr == "" {
		return query, nil
	}

	// TODO(patrik): Better errors
	expr, err := fullParseFilter(resolver, filterStr)
	if err != nil {
		return nil, err
	}

	return query.Where(expr), nil
}

func applySort(query *goqu.SelectDataset, resolver *filter.Resolver, sortStr string) (*goqu.SelectDataset, error) {
	sortExpr, err := filter.ParseSort(sortStr)
	if err != nil {
		return nil, InvalidSort(err)
	}

	sortExpr, err = resolver.ResolveSort(sortExpr)
	if err != nil {
		return nil, InvalidSort(err)
	}

	exprs, err := generateSort(sortExpr)
	if err != nil {
		return nil, InvalidSort(err)
	}

	return query.Order(exprs...), nil
}

func fullParseFilter(resolver *filter.Resolver, filterStr string) (exp.Expression, error) {
	ast, err := parser.ParseExpr(filterStr)
	if err != nil {
		return nil, InvalidFilter(err)
	}

	e, err := resolver.Resolve(ast)
	if err != nil {
		return nil, InvalidFilter(err)
	}

	re, err := generateFilter(e)
	if err != nil {
		return nil, InvalidFilter(err)
	}

	return re, nil
}

func generateTableSelect(table *filter.Table, ids []string) *goqu.SelectDataset {
	return goqu.From(table.Name).
		Select(table.SelectName).
		Where(goqu.I(table.WhereName).In(ids))
}

var opMapping = map[filter.OpKind]string{
	filter.OpEqual:        "(? == ?)",
	filter.OpNotEqual:     "(? != ?)",
	filter.OpLike:         "(? LIKE ?)",
	filter.OpGreater:      "(? > ?)",
	filter.OpGreaterEqual: "(? >= ?)",
	filter.OpLesser:       "(? < ?)",
	filter.OpLesserEqual:  "(? <= ?)",
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
	case *filter.IsNullExpr:
		if e.Not {
			return goqu.L("(? IS NOT NULL)", goqu.I(e.Name)), nil
		} else {
			return goqu.L("(? IS NULL)", goqu.I(e.Name)), nil
		}
	case *filter.OpExpr:
		if op, ok := opMapping[e.Kind]; ok {
			return goqu.L(op, goqu.I(e.Name), e.Value), nil
		}

		return nil, fmt.Errorf("unimplemented OpKind %d", e.Kind)
	case *filter.InTableExpr:
		s := generateTableSelect(&e.Table, e.Ids)

		if e.Not {
			// TODO(patrik): Move tracks.id
			return goqu.L("? NOT IN ?", goqu.I(e.IdSelector), s), nil
		} else {
			return goqu.L("? IN ?", goqu.I(e.IdSelector), s), nil
		}
	}

	return nil, fmt.Errorf("Unimplemented expr %T", e)
}

func generateSort(e filter.SortExpr) ([]exp.OrderedExpression, error) {
	switch e := e.(type) {
	case *filter.SortExprSort:
		var items []exp.OrderedExpression
		for _, item := range e.Items {
			switch item.Type {
			case filter.SortTypeAsc:
				items = append(items, goqu.I(item.Name).Asc().NullsLast())
			case filter.SortTypeDesc:
				items = append(items, goqu.I(item.Name).Desc().NullsFirst())
			default:
				return nil, fmt.Errorf("Unknown SortItemType: %d", item.Type)
			}
		}

		return items, nil
	case *filter.SortExprRandom:
		oe := goqu.Func("RANDOM").Asc()
		return []exp.OrderedExpression{oe}, nil
	}

	return nil, fmt.Errorf("Unimplemented expr %T", e)
}
