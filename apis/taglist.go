package apis

import (
	"context"
	"errors"
	"go/parser"
	"net/http"
	"strconv"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/database/adapter"
	"github.com/nanoteck137/dwebble/tools/filter"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/nanoteck137/validate"
)

type Taglist struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Filter string `json:"filter"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

type GetTaglists struct {
	Taglists []Taglist `json:"taglists"`
}

type GetTaglistById struct {
	Taglist
}

type GetTaglistTracks struct {
	Page   types.Page `json:"page"`
	Tracks []Track    `json:"tracks"`
}

type CreateTaglist struct {
	Id string `json:"id"`
}

type CreateTaglistBody struct {
	Name   string `json:"name"`
	Filter string `json:"filter"`
}

func (b *CreateTaglistBody) Transform() {
	b.Name = transform.String(b.Name)
	b.Filter = transform.String(b.Filter)
}

func (b CreateTaglistBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
		validate.Field(&b.Filter, validate.Required),
	)
}

// TODO(patrik): Validation and transform
type UpdateTaglistBody struct {
	Name   *string `json:"name,omitempty"`
	Filter *string `json:"filter,omitempty"`
}

func (b *UpdateTaglistBody) Transform() {
	b.Name = transform.StringPtr(b.Name)
	b.Filter = transform.StringPtr(b.Filter)
}

func (b UpdateTaglistBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),
		validate.Field(&b.Filter, validate.Required.When(b.Filter != nil)),
	)
}

func TestFilter(filterStr string) error {
	ast, err := parser.ParseExpr(filterStr)
	if err != nil {
		return InvalidFilter(err)
	}

	a := adapter.TrackResolverAdapter{}
	r := filter.New(&a)
	_, err = r.Resolve(ast)
	if err != nil {
		return InvalidFilter(err)
	}

	return nil
}

func InstallTaglistHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetTaglists",
			Path:         "/taglists",
			Method:       http.MethodGet,
			ResponseType: GetTaglists{},
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
			Name:         "GetTaglistById",
			Path:         "/taglists/:id",
			Method:       http.MethodGet,
			ResponseType: GetTaglistById{},
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
			Name:         "GetTaglistTracks",
			Path:         "/taglists/:id/tracks",
			Method:       http.MethodGet,
			ResponseType: GetTaglistTracks{},
			Errors:       []pyrin.ErrorType{ErrTypeTaglistNotFound},
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
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, TaglistNotFound()
					}

					return nil, err
				}

				if taglist.OwnerId != user.Id {
					return nil, TaglistNotFound()
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
					Tracks: make([]Track, len(tracks)),
				}

				for i, track := range tracks {
					res.Tracks[i] = ConvertDBTrack(c, track)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreateTaglist",
			Path:         "/taglists",
			Method:       http.MethodPost,
			ResponseType: CreateTaglist{},
			BodyType:     CreateTaglistBody{},
			Errors:       []pyrin.ErrorType{ErrTypeInvalidFilter},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[CreateTaglistBody](c)
				if err != nil {
					return nil, err
				}

				err = TestFilter(body.Filter)
				if err != nil {
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
			Name:     "DeleteTaglist",
			Path:     "/taglists/:id",
			Method:   http.MethodPatch,
			BodyType: UpdateTaglistBody{},
			Errors:   []pyrin.ErrorType{ErrTypeTaglistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				taglistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				taglist, err := app.DB().GetTaglistById(c.Request().Context(), taglistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, TaglistNotFound()
					}

					return nil, err
				}

				if taglist.OwnerId != user.Id {
					return nil, TaglistNotFound()
				}

				ctx := context.TODO()

				err = app.DB().RemoveTaglist(ctx, taglist.Id)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "UpdateTaglist",
			Path:     "/taglists/:id",
			Method:   http.MethodPatch,
			BodyType: UpdateTaglistBody{},
			Errors:   []pyrin.ErrorType{ErrTypeTaglistNotFound, ErrTypeInvalidFilter},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				taglistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[UpdateTaglistBody](c)
				if err != nil {
					return nil, err
				}

				taglist, err := app.DB().GetTaglistById(c.Request().Context(), taglistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, TaglistNotFound()
					}

					return nil, err
				}

				if taglist.OwnerId != user.Id {
					return nil, TaglistNotFound()
				}

				changes := database.TaglistChanges{}

				if body.Name != nil {
					changes.Name = types.Change[string]{
						Value:   *body.Name,
						Changed: *body.Name != taglist.Name,
					}
				}

				if body.Filter != nil {
					changes.Filter = types.Change[string]{
						Value:   *body.Filter,
						Changed: *body.Filter != taglist.Filter,
					}
				}

				if changes.Filter.Changed {
					err = TestFilter(changes.Filter.Value)
					if err != nil {
						return nil, err
					}
				}

				err = app.DB().UpdateTaglist(c.Request().Context(), taglist.Id, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
