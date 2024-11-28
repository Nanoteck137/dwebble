package apis

import (
	"context"
	"errors"
	"go/parser"
	"net/http"
	"strconv"
	"strings"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/filter"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/validate"
)

type Taglist struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Filter  string `json:"filter"`

	Created int64  `json:"created"`
	Updated int64  `json:"updated"`
}

type GetTaglists struct {
	Taglists []Taglist `json:"taglists"`
}

type GetTaglistById struct {
	Taglist
}

type GetTaglistTracks struct {
	Page   types.Page    `json:"page"`
	Tracks []types.Track `json:"tracks"`
}

type CreateTaglist struct {
	Id string `json:"id"`
}

var _ pyrin.Body = (*CreateTaglistBody)(nil)

type CreateTaglistBody struct {
	Name   string `json:"name"`
	Filter string `json:"filter"`
}

func (b CreateTaglistBody) Validate(validator validate.Validator) error {
	panic("unimplemented")
}

var _ pyrin.Body = (*UpdateTaglistBody)(nil)

type UpdateTaglistBody struct {
	Name   *string `json:"name,omitempty"`
	Filter *string `json:"filter,omitempty"`
}

func (b UpdateTaglistBody) Validate(validator validate.Validator) error {
	panic("unimplemented")
}

func InstallTaglistHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:     "GetTaglists",
			Path:     "/taglists",
			Method:   http.MethodGet,
			DataType: GetTaglists{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				taglists, err := app.DB().GetTaglistByUser(c.Request().Context(), user.Id)
				if err != nil {
					return nil, err
				}

				res := GetTaglists{
					Taglists: make([]Taglist, len(taglists)),
				}

				for i, taglist := range taglists {
					res.Taglists[i] = Taglist{
						Id:      taglist.Id,
						Name:    taglist.Name,
						Filter:  taglist.Filter,
						Created: taglist.Created,
						Updated: taglist.Updated,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "GetTaglistById",
			Path:     "/taglists/:id",
			Method:   http.MethodGet,
			DataType: GetTaglistById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				taglistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				taglist, err := app.DB().GetTaglistById(c.Request().Context(), taglistId)
				if err != nil {
					return nil, err
				}

				if taglist.OwnerId != user.Id {
					// TODO(patrik): Fix error
					return nil, errors.New("No taglist")
				}

				res := GetTaglistById{
					Taglist: Taglist{
						Id:      taglist.Id,
						Name:    taglist.Name,
						Filter:  taglist.Filter,
						Created: taglist.Created,
						Updated: taglist.Updated,
					},
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "GetTaglistTracks",
			Path:     "/taglists/:id/tracks",
			Method:   http.MethodGet,
			DataType: GetTaglistTracks{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				taglistId := c.Param("id")

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				taglist, err := app.DB().GetTaglistById(ctx, taglistId)
				if err != nil {
					return nil, err
				}

				if taglist.OwnerId != user.Id {
					// TODO(patrik): Fix error
					return nil, errors.New("No taglist")
				}

				// TODO(patrik): Standardize FetchOption
				perPage := 100
				page := 0

				if s := q.Get("perPage"); s != "" {
					i, err := strconv.Atoi(s)
					if err != nil {
						return nil, err
					}

					perPage = i
				}

				if s := q.Get("page"); s != "" {
					i, err := strconv.Atoi(s)
					if err != nil {
						return nil, err
					}

					page = i
				}

				opts := database.FetchOption{
					Filter:  taglist.Filter,
					Sort:    q.Get("sort"),
					PerPage: perPage,
					Page:    page,
				}

				tracks, p, err := app.DB().GetPagedTracks(ctx, opts)
				if err != nil {
					if errors.Is(err, database.ErrInvalidFilter) {
						return nil, InvalidFilter(err)
					}

					if errors.Is(err, database.ErrInvalidSort) {
						return nil, InvalidSort(err)
					}

					return nil, err
				}

				res := GetTaglistTracks{
					Page:   p,
					Tracks: make([]types.Track, len(tracks)),
				}

				for i, track := range tracks {
					res.Tracks[i] = ConvertDBTrack(c, track)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "CreateTaglist",
			Path:     "/taglists",
			Method:   http.MethodPost,
			DataType: CreateTaglist{},
			BodyType: CreateTaglistBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := Body[CreateTaglistBody](c)
				if err != nil {
					return nil, err
				}

				ast, err := parser.ParseExpr(body.Filter)
				if err != nil {
					// TODO(patrik): Better errors
					return nil, err
				}

				a := database.TrackResolverAdapter{}
				r := filter.New(&a)
				_, err = r.Resolve(ast)
				if err != nil {
					// TODO(patrik): Better errors
					return nil, err
				}

				res, err := app.DB().CreateTaglist(c.Request().Context(), database.CreateTaglistParams{
					Name:    body.Name,
					Filter:  body.Filter,
					OwnerId: user.Id,
				})
				if err != nil {
					return nil, err
				}

				return CreateTaglist{
					Id: res.Id,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "UpdateTaglist",
			Path:     "/taglists/:id",
			Method:   http.MethodPatch,
			BodyType: UpdateTaglistBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				taglistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := Body[UpdateTaglistBody](c)
				if err != nil {
					return nil, err
				}

				taglist, err := app.DB().GetTaglistById(c.Request().Context(), taglistId)
				if err != nil {
					return nil, err
				}

				if taglist.OwnerId != user.Id {
					// TODO(patrik): Fix error
					return nil, errors.New("No taglist")
				}

				var name types.Change[string]
				if body.Name != nil {
					n := strings.TrimSpace(*body.Name)
					name.Value = n
					name.Changed = n != taglist.Name
				}

				var filterChange types.Change[string]
				if body.Filter != nil {
					f := strings.TrimSpace(*body.Filter)
					filterChange.Value = f
					filterChange.Changed = f != taglist.Filter
				}

				if filterChange.Changed {
					ast, err := parser.ParseExpr(filterChange.Value)
					if err != nil {
						// TODO(patrik): Better errors
						return nil, err
					}

					a := database.TrackResolverAdapter{}
					r := filter.New(&a)
					_, err = r.Resolve(ast)
					if err != nil {
						// TODO(patrik): Better errors
						return nil, err
					}
				}

				err = app.DB().UpdateTaglist(c.Request().Context(), taglist.Id, database.TaglistChanges{
					Name:   name,
					Filter: filterChange,
				})
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
