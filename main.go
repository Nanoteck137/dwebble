package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/nanoteck137/dwebble/api"

	_ "github.com/nanoteck137/dwebble/docs"
)

type WorkDir string

func (d WorkDir) TracksDir() string {
	dir := string(d)
	return path.Join(dir, "tracks")
}

func (d WorkDir) ImagesDir() string {
	dir := string(d)
	return path.Join(dir, "images")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL not set")
	}

	db, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 1 * 1024 * 1024 * 1024,
	})

	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	workDirPath := os.Getenv("WORK_DIR")
	if workDirPath == "" {
		log.Println("Warning: WORK_DIR not set, using cwd")
		workDirPath = "./"
	}
	workDir := WorkDir(workDirPath)

	app.Use("/tracks", filesystem.New(filesystem.Config{
		Root: http.Dir(workDir.TracksDir()),
	}))

	app.Use("/images", filesystem.New(filesystem.Config{
		Root: http.Dir(workDir.ImagesDir()),
	}))

	app.Mount("/api/v1", api.New(db))

	log.Fatal(app.Listen(":3000"))
}
