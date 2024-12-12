package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/tools/utils"
)

type ApiToken struct {
	Id     string `db:"id"`
	UserId string `db:"user_id"`

	Name string `db:"name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func ApiTokenQuery() *goqu.SelectDataset {
	query := dialect.From("api_tokens").
		Select(
			"api_tokens.id",
			"api_tokens.user_id",

			"api_tokens.name",

			"api_tokens.updated",
			"api_tokens.created",
		).
		Prepared(true)

	return query
}

func (db *Database) GetApiTokenById(ctx context.Context, id string) (ApiToken, error) {
	query := ApiTokenQuery().
		Where(goqu.I("api_tokens.id").Eq(id))

	var item ApiToken
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ApiToken{}, ErrItemNotFound
		}

		return ApiToken{}, err
	}

	return item, nil
}

func (db *Database) GetAllApiTokensForUser(ctx context.Context, userId string) ([]ApiToken, error) {
	query := ApiTokenQuery().
		Where(goqu.I("api_tokens.user_id").Eq(userId))

	var items []ApiToken
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

type CreateApiTokenParams struct {
	Id     string
	UserId string
	Name   string

	Created int64
	Updated int64
}

func (db *Database) CreateApiToken(ctx context.Context, params CreateApiTokenParams) (ApiToken, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateApiTokenId()
	}

	query := dialect.Insert("api_tokens").Rows(goqu.Record{
		"id":      id,
		"user_id": params.UserId,

		"name": params.Name,

		"created": created,
		"updated": updated,
	}).
		Returning(
			"api_tokens.id",
			"api_tokens.user_id",

			"api_tokens.name",

			"api_tokens.updated",
			"api_tokens.created",
		).
		Prepared(true)

	var item ApiToken
	err := db.Get(&item, query)
	if err != nil {
		return ApiToken{}, err
	}

	return item, nil
}

func (db *Database) RemoveApiToken(ctx context.Context, id string) error {
	query := dialect.Delete("api_tokens").
		Where(goqu.I("api_tokens.id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
