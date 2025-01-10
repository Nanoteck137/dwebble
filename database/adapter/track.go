package adapter

import (
	"go/ast"

	"github.com/nanoteck137/dwebble/tools/filter"
	"github.com/nanoteck137/dwebble/tools/utils"
)

var _ filter.ResolverAdapter = (*TrackResolverAdapter)(nil)

type TrackResolverAdapter struct{}

func (a *TrackResolverAdapter) DefaultSort() (string, filter.SortType) {
	return "tracks.name", filter.SortTypeAsc
}

func (a *TrackResolverAdapter) ResolveVariableName(name string) (filter.Name, bool) {
	switch name {
	case "id":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "tracks.id",
		}, true
	case "name":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "tracks.name",
		}, true
	case "otherName":
		return filter.Name{
			Kind:     filter.NameKindString,
			Name:     "tracks.other_name",
			Nullable: true,
		}, true
	case "number":
		return filter.Name{
			Kind:     filter.NameKindNumber,
			Name:     "tracks.track_number",
			Nullable: true,
		}, true
	case "duration":
		return filter.Name{
			Kind:     filter.NameKindNumber,
			Name:     "tracks.duration",
			Nullable: true,
		}, true
	case "year":
		return filter.Name{
			Kind:     filter.NameKindNumber,
			Name:     "tracks.year",
			Nullable: true,
		}, true
	case "albumId":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "tracks.album_id",
		}, true
	case "artistId":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "tracks.artist_id",
		}, true
	case "albumName":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "albums.name",
		}, true
	case "albumOtherName":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "albums.other_name",
			Nullable: true,
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
			Name: "tracks.created",
		}, true
	case "updated":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "tracks.updated",
		}, true
	}

	return filter.Name{}, false
}

func (a *TrackResolverAdapter) ResolveNameToId(typ, name string) (string, bool) {
	switch typ {
	case "tags":
		return utils.Slug(name), true
	case "featuringArtists":
		return name, true
	}

	return "", false
}

func (a *TrackResolverAdapter) ResolveTable(typ string) (filter.Table, bool) {
	switch typ {
	case "tags":
		return filter.Table{
			Name:       "tracks_tags",
			SelectName: "track_id",
			WhereName:  "tag_slug",
		}, true
	case "featuringArtists":
		return filter.Table{
			Name:       "tracks_featuring_artists",
			SelectName: "track_id",
			WhereName:  "artist_id",
		}, true
	}

	return filter.Table{}, false
}

func (a *TrackResolverAdapter) ResolveFunctionCall(resolver *filter.Resolver, name string, args []ast.Expr) (filter.FilterExpr, error) {
	switch name {
	case "hasTag":
		return resolver.InTable(name, "tags", "tracks.id", args)
	case "hasFeaturingArtist":
		return resolver.InTable(name, "featuringArtists", "tracks.id", args)
	}

	return nil, filter.UnknownFunction(name)
}
