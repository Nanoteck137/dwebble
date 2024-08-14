package database

import (
	"context"
	"database/sql"
	"fmt"
	"go/parser"

	"github.com/doug-martin/goqu/v9"
	goqusqlite3 "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/nanoteck137/dwebble/tools/filter"
	"github.com/nanoteck137/dwebble/types"
	"github.com/pressly/goose/v3"

	_ "github.com/mattn/go-sqlite3"
)

var dialect = goqu.Dialect("sqlite_returning")

type ToSQL interface {
	ToSQL() (string, []interface{}, error)
}

type Connection interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

type Database struct {
	RawConn *sql.DB
	Conn    Connection
}

func New(conn *sql.DB) *Database {
	return &Database{
		RawConn: conn,
		Conn:    conn,
	}
}

func Open(workDir types.WorkDir) (*Database, error) {
	dbUrl := fmt.Sprintf("file:%s?_foreign_keys=true", workDir.DatabaseFile())

	conn, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		return nil, err
	}

	return New(conn), nil
}

func (db *Database) Begin() (*Database, *sql.Tx, error) {
	tx, err := db.RawConn.Begin()
	if err != nil {
		return nil, nil, err
	}

	return &Database{
		RawConn: db.RawConn,
		Conn:    tx,
	}, tx, nil
}

func (db *Database) Query(ctx context.Context, s ToSQL) (*sql.Rows, error) {
	sql, params, err := s.ToSQL()
	if err != nil {
		return nil, err
	}

	return db.Conn.Query(sql, params...)
}

func (db *Database) QueryRow(ctx context.Context, s ToSQL) (*sql.Row, error) {
	sql, params, err := s.ToSQL()
	if err != nil {
		return nil, err
	}

	row := db.Conn.QueryRow(sql, params...)
	return row, nil
}

func (db *Database) Exec(ctx context.Context, s ToSQL) (sql.Result, error) {
	sql, params, err := s.ToSQL()
	if err != nil {
		return nil, err
	}

	return db.Conn.Exec(sql, params...)
}

func init() {
	opts := goqusqlite3.DialectOptions()
	opts.SupportsReturn = true
	goqu.RegisterDialect("sqlite_returning", opts)
}

func fullParseFilter(adapter filter.ResolverAdapter, filterStr string) (exp.Expression, error) {
	ast, err := parser.ParseExpr(filterStr)
	if err != nil {
		return nil, types.NewApiError(400, "Failed to parse filter: " + err.Error())
	}

	r := filter.New(adapter)
	e, err := r.Resolve(ast)
	if err != nil {
		return nil, types.NewApiError(400, err.Error())
	}

	re, err := generateFilter(e)
	if err != nil {
		return nil, err
	}

	return re, nil
}

func RunMigrateUp(db *Database) error {
	return goose.Up(db.RawConn, ".")
}
