package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/dwebble/v2/internal/database"
	"github.com/dwebble/v2/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

//go:embed images/*
var content embed.FS

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

	queries := database.New(db)

	app := fiber.New()
	app.Use(logger.New())

	app.Use("/images", filesystem.New(filesystem.Config{
		Root: http.Dir("./images"),
	}))

	app.Use("/images", filesystem.New(filesystem.Config{
		Root: http.FS(content),
		PathPrefix: "images",
	}))

	apiConfig := handlers.New(queries)

	v1 := app.Group("/api/v1")

	v1.Get("/artists", apiConfig.HandlerGetAllArtists)
	v1.Get("/artists/:id", apiConfig.HandlerGetArtist)

	v1.Get("/albums", apiConfig.HandlerGetAllAlbums)
	v1.Get("/albums/:id", apiConfig.HandlerGetAlbum)

	log.Fatal(app.Listen(":3000"))
}
