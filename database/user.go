package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type User struct {
	Id       string `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

func UserQuery() *goqu.SelectDataset {
	query := dialect.From("users").
		Select(
			"users.id", 
			"users.username",
			"users.password",
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

func (db *Database) CreateUser(ctx context.Context, username, password string) (User, error) {
	query := dialect.
		Insert("users").
		Rows(goqu.Record{
			"id":       utils.CreateId(),
			"username": username,
			"password": password,
		}).
		Returning(
			"users.id", 
			"users.username", 
			"users.password",
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
}

func (db *Database) UpdateUser(ctx context.Context, id string, changes UserChanges) error {
	record := goqu.Record{}

	addToRecord(record, "username", changes.Username)
	addToRecord(record, "password", changes.Password)

	if len(record) == 0 {
		return nil
	}

	// record["updated"] = time.Now().UnixMilli()

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
