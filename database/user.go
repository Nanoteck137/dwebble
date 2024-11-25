package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type User struct {
	Id       string
	Username string
	Password string
}

func (db *Database) GetAllUsers(ctx context.Context) ([]User, error) {
	ds := dialect.From("users").
		Select("id", "username")

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []User
	for rows.Next() {
		var item User
		err := rows.Scan(&item.Id, &item.Username)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (db *Database) CreateUser(ctx context.Context, username, password string) (User, error) {
	ds := dialect.
		Insert("users").
		Rows(goqu.Record{
			"id":       utils.CreateId(),
			"username": username,
			"password": password,
		}).
		Returning("id", "username", "password")

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return User{}, err
	}

	var item User
	err = row.Scan(&item.Id, &item.Username, &item.Password)
	if err != nil {
		return User{}, err
	}

	return item, nil
}

func (db *Database) GetUserById(ctx context.Context, id string) (User, error) {
	ds := dialect.
		From("users").
		Select("id", "username", "password").
		Where(goqu.C("id").Eq(id))

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return User{}, err
	}

	var item User
	err = row.Scan(&item.Id, &item.Username, &item.Password)
	if err != nil {
		return User{}, err
	}

	return item, nil
}

func (db *Database) GetUserByUsername(ctx context.Context, username string) (User, error) {
	ds := dialect.
		From("users").
		Select("id", "username", "password").
		Where(goqu.C("username").Eq(username))
		// TODO(patrik): Prepared

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return User{}, err
	}

	var item User
	err = row.Scan(&item.Id, &item.Username, &item.Password)
	if err != nil {
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
		Where(goqu.I("id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}
