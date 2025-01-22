package filter

import (
	"errors"
	"fmt"
	"strings"
)

// TODO(patrik): Rename to SortOrder
type SortType int

const (
	SortTypeAsc SortType = iota
	SortTypeDesc
)

type SortExpr interface {
	sortType()
}

type SortItem struct {
	Type SortType
	Name string
}

type SortExprSort struct {
	Items []SortItem
}

type SortExprRandom struct{}
type SortExprDefault struct{}

func (e *SortExprSort) sortType()    {}
func (e *SortExprRandom) sortType()  {}
func (e *SortExprDefault) sortType() {}

func ParseSort(s string) (SortExpr, error) {
	// TODO(patrik): Trim the strings
	split := strings.Split(s, "=")

	mode := split[0]
	switch mode {
	case "sort":
		args := strings.Split(split[1], ",")

		var items []SortItem

		for _, arg := range args {
			var item SortItem
			switch arg[0] {
			case '+':
				item = SortItem{
					Type: SortTypeAsc,
					Name: arg[1:],
				}
			case '-':
				item = SortItem{
					Type: SortTypeDesc,
					Name: arg[1:],
				}
			default:
				item = SortItem{
					Type: SortTypeAsc,
					Name: arg,
				}
			}

			items = append(items, item)
		}

		return &SortExprSort{
			Items: items,
		}, nil
	case "random":
		return &SortExprRandom{}, nil
	case "", "default":
		return &SortExprDefault{}, nil
	default:
		return nil, errors.New("unknown sort mode")
	}
}

func (r *Resolver) ResolveSort(e SortExpr) (SortExpr, error) {
	switch e := e.(type) {
	case *SortExprSort:
		for i, item := range e.Items {
			resolvedName, ok := r.adapter.ResolveVariableName(item.Name)
			if !ok {
				return nil, UnknownName(item.Name)
			}

			e.Items[i].Name = resolvedName.Name
		}

		return e, nil
	case *SortExprRandom:
		return e, nil
	case *SortExprDefault:
		defaultSort, typ := r.adapter.DefaultSort()
		return &SortExprSort{
			Items: []SortItem{
				{
					Type: typ,
					Name: defaultSort,
				},
			},
		}, nil
	}

	// TODO(patrik): Internal error
	return nil, fmt.Errorf("unimplemented expr %T", e)
}
