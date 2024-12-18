package adapter

import (
	"go/ast"

	"github.com/nanoteck137/dwebble/tools/filter"
)

var _ filter.ResolverAdapter = (*TrackResolverAdapter)(nil)

type AlbumResolverAdapter struct{}

func (a *AlbumResolverAdapter) DefaultSort() (string, filter.SortType) {
	return "albums.name", filter.SortTypeAsc
}

func (a *AlbumResolverAdapter) ResolveVariableName(name string) (filter.Name, bool) {
	switch name {
	case "id":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "albums.id",
		}, true
	case "name":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "albums.name",
		}, true
	case "otherName":
		return filter.Name{
			Kind:     filter.NameKindString,
			Name:     "albums.other_name",
			Nullable: true,
		}, true
	case "year":
		return filter.Name{
			Kind:     filter.NameKindNumber,
			Name:     "albums.year",
			Nullable: true,
		}, true
	case "artistId":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "albums.artist_id",
		}, true
	case "artistName":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "artists.name",
		}, true
	case "artistOtherName":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "artists.other_name",
			Nullable: true,
		}, true
	case "created":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "albums.created",
		}, true
	case "updated":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "albums.updated",
		}, true
	}

	return filter.Name{}, false
}

func (a *AlbumResolverAdapter) ResolveNameToId(typ, name string) (string, bool) {
	// switch typ {
	// case "tags":
	// 	return utils.Slug(name), true
	// }

	return "", false
}

func (a *AlbumResolverAdapter) ResolveTable(typ string) (filter.Table, bool) {
	// switch typ {
	// case "tags":
	// 	return filter.Table{
	// 		Name:       "tracks_to_tags",
	// 		SelectName: "track_id",
	// 		WhereName:  "tag_slug",
	// 	}, true
	// }

	return filter.Table{}, false
}

func (a *AlbumResolverAdapter) ResolveFunctionCall(resolver *filter.Resolver, name string, args []ast.Expr) (filter.FilterExpr, error) {
	// switch name {
	// case "hasTag":
	// 	return resolver.InTable(name, "tags", args)
	// }

	return nil, filter.UnknownFunction(name)
}
