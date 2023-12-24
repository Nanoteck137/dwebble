package main

import (
	"context"
	"log"
	"os"

	"github.com/dwebble/v2/internal/database"
	"github.com/dwebble/v2/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

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

	queries.CreateArtist(context.Background(), database.CreateArtistParams{ID: "hello", Name: "Text"})

	apiConfig := handlers.New(queries)

	v1 := app.Group("/api/v1")

	v1.Get("/artists", apiConfig.HandlerGetAllArtists)
	v1.Get("/artists/:id", apiConfig.HandlerGetArtist)

	v1.Get("/albums", apiConfig.HandlerGetAllAlbums)
	v1.Get("/albums/:id", apiConfig.HandlerGetAlbum)

	log.Fatal(app.Listen(":3000"))
}
