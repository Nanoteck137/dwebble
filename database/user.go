package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type UserSettings struct {
	Id            string         `db:"id"`
	DisplayName   sql.NullString `db:"display_name"`
	QuickPlaylist sql.NullString `db:"quick_playlist"`
}

type User struct {
	Id       string `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Role     string `db:"role"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	// NOTE(patrik): This needs to match UserSettings
	DisplayName   sql.NullString `db:"display_name"`
	QuickPlaylist sql.NullString `db:"quick_playlist"`
}

func (u User) ToUserSettings() UserSettings {
	return UserSettings{
		Id:            u.Id,
		DisplayName:   u.DisplayName,
		QuickPlaylist: u.QuickPlaylist,
	}
}

func UserQuery() *goqu.SelectDataset {
	query := dialect.From("users").
		Select(
			"users.id",
			"users.username",
			"users.password",
			"users.role",

			"users.created",
			"users.updated",

			"users_settings.display_name",
			"users_settings.quick_playlist",
		).
		Prepared(true).
		LeftJoin(
			goqu.I("users_settings"),
			goqu.On(goqu.I("users.id").Eq(goqu.I("users_settings.id"))),
		)

	return query
}

func UserSettingsQuery() *goqu.SelectDataset {
	query := dialect.From("users_settings").
		Select(
			"users_settings.id",
			"users_settings.display_name",

			"users_settings.quick_playlist",

			"users_settings.created",
			"users_settings.updated",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllUsers(ctx context.Context) ([]User, error) {
	query := UserQuery()

	var items []User
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

type CreateUserParams struct {
	Id       string
	Username string
	Password string
	Role     string

	Created int64
	Updated int64
}

func (db *Database) CreateUser(ctx context.Context, params CreateUserParams) (User, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateId()
	}

	query := dialect.
		Insert("users").
		Rows(goqu.Record{
			"id":       utils.CreateId(),
			"username": params.Username,
			"password": params.Password,
			"role":     params.Role,

			"created": created,
			"updated": updated,
		}).
		Returning(
			"users.id",
			"users.username",
			"users.password",
			"users.role",

			"users.created",
			"users.updated",
		)

	var item User
	err := db.Get(&item, query)
	if err != nil {
		return User{}, err
	}

	return item, nil
}

func (db *Database) GetUserById(ctx context.Context, id string) (User, error) {
	query := UserQuery().
		Where(goqu.I("users.id").Eq(id))

	var item User
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrItemNotFound
		}

		return User{}, err
	}

	return item, nil
}

func (db *Database) GetUserSettingsById(ctx context.Context, id string) (UserSettings, error) {
	query := UserSettingsQuery().
		Where(goqu.I("users_settings.id").Eq(id))

	var item UserSettings
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserSettings{}, ErrItemNotFound
		}

		return UserSettings{}, err
	}

	return item, nil
}

func (db *Database) GetUserByUsername(ctx context.Context, username string) (User, error) {
	query := UserQuery().
		Where(goqu.I("users.username").Eq(username))

	var item User
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrItemNotFound
		}

		return User{}, err
	}

	return item, nil
}

type UserChanges struct {
	Username types.Change[string]
	Password types.Change[string]

	Created types.Change[int64]
}

func (db *Database) UpdateUser(ctx context.Context, id string, changes UserChanges) error {
	record := goqu.Record{}

	addToRecord(record, "username", changes.Username)
	addToRecord(record, "password", changes.Password)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("users").
		Set(record).
		Prepared(true).
		Where(goqu.I("users.id").Eq(id))

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateUserSettings(ctx context.Context, settings UserSettings) error {
	query := dialect.Insert("users_settings").
		Rows(goqu.Record{
			"id":             settings.Id,
			"display_name":   settings.DisplayName,
			"quick_playlist": settings.QuickPlaylist,
		}).
		Prepared(true).
		OnConflict(goqu.DoUpdate("id", goqu.Record{
			"display_name":   settings.DisplayName,
			"quick_playlist": settings.QuickPlaylist,
		}))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
