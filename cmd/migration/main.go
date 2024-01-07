package main

import (
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/urfave/cli/v2"
)

const sqlDir = "./sql/schema"

func main() {
	godotenv.Load()

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL not set")
	}

	db, err := goose.OpenDBWithDriver("pgx", dbUrl)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	app := cli.App{
		Usage: "Hello World",

		Commands: []*cli.Command{
			{
				Name: "up",
				Action: func(ctx *cli.Context) error {
					return goose.Up(db, sqlDir)
				},
			},
			{
				Name: "down",
				Action: func(ctx *cli.Context) error {
					return goose.Down(db, sqlDir)
				},
			},
			{
				Name: "reset",
				Action: func(ctx *cli.Context) error {
					return goose.Reset(db, sqlDir)
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
