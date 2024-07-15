package sort

import (
	"fmt"
	"strings"

	"github.com/nanoteck137/dwebble/sortfilter/filter"
)

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

func Parse(s string) (SortExpr, error) {
	split := strings.Split(s, "=")
	fmt.Printf("split: %v\n", split)

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
	}

	return &SortExprDefault{}, nil
}

type Resolver struct {
	adapter filter.ResolverAdapter
}

func New(adapter filter.ResolverAdapter) *Resolver {
	return &Resolver{
		adapter: adapter,
	}
}

func (r *Resolver) Resolve(e SortExpr) (SortExpr, error) {
	switch e := e.(type) {
	case *SortExprSort:
		for i, item := range e.Items {
			resolvedName, err := r.adapter.MapName(item.Name)
			if err != nil {
				return nil, err
			}

			e.Items[i].Name = resolvedName.Name
		}

		return e, nil
	case *SortExprRandom:
		return e, nil
	case *SortExprDefault:
		defaultSort := r.adapter.GetDefaultSort()
		return &SortExprSort{
			Items: []SortItem{
				{
					Type: SortTypeAsc,
					Name: defaultSort,
				},
			},
		}, nil
	}

	return nil, fmt.Errorf("Unimplemented expr %T", e)
}
