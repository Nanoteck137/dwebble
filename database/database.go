package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ToSQL interface {
	ToSQL() (string, []interface{}, error)
}

type Database struct {
	conn *pgxpool.Pool
}

func (db *Database) Query(ctx context.Context, s ToSQL) (pgx.Rows, error) {
	sql, params, err := s.ToSQL()
	if err != nil {
		return nil, err
	}

	return db.conn.Query(ctx, sql, params...)
}

func (db *Database) QueryRow(ctx context.Context, s ToSQL) (pgx.Row, error) {
	sql, params, err := s.ToSQL()
	if err != nil {
		return nil, err
	}

	row := db.conn.QueryRow(ctx, sql, params...)
	return row, nil
}

func (db *Database) Exec(ctx context.Context, s ToSQL) (pgconn.CommandTag, error) {
	sql, params, err := s.ToSQL()
	if err != nil {
		return pgconn.NewCommandTag(""), err
	}

	return db.conn.Exec(ctx, sql, params...)
}

var dialect = goqu.Dialect("postgres")

func New(conn *pgxpool.Pool) *Database {
	return &Database{
		conn: conn,
	}
}

