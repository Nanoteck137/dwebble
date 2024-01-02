package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/nanoteck137/dwebble/v2/internal/database"
	"github.com/nanoteck137/dwebble/v2/internal/handlers"
	"github.com/nanoteck137/dwebble/v2/views"
)

//go:embed images/* public/*
var content embed.FS

func component(c *fiber.Ctx, component templ.Component) error {
	c.Response().Header.SetContentType(fiber.MIMETextHTML)
	return component.Render(c.Context(), c.Response().BodyWriter())
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

	queries := database.New(db)

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(logger.New())

	app.Use("/tracks", filesystem.New(filesystem.Config{
		Root: http.Dir("./test"),
	}))

	app.Use("/images", filesystem.New(filesystem.Config{
		Root: http.Dir("./images"),
	}))

	app.Use("/images", filesystem.New(filesystem.Config{
		Root:       http.FS(content),
		PathPrefix: "images",
	}))

	app.Use("/public", filesystem.New(filesystem.Config{
		Root:       http.FS(content),
		PathPrefix: "public",
	}))

	apiConfig := handlers.New(queries)

	app.Get("/", func(c *fiber.Ctx) error {
		artists, err := queries.GetAllArtists(c.Context())
		if err != nil {
			return err
		}

		return component(c, views.Index(artists))
	})

	app.Get("/artist/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		artist, err := queries.GetArtist(c.Context(), id)
		if err != nil {
			return err
		}

		return component(c, views.ArtistPage(artist))
	})

	v1 := app.Group("/api/v1")

	v1.Get("/artists", apiConfig.HandlerGetAllArtists)
	v1.Get("/artists/:id", apiConfig.HandlerGetArtist)

	v1.Get("/albums", apiConfig.HandlerGetAllAlbums)
	v1.Get("/albums/:id", apiConfig.HandlerGetAlbum)

	v1.Get("/tracks", apiConfig.HandlerGetAllTracks)
	v1.Get("/tracks/:id", apiConfig.HandlerGetTrack)

	log.Fatal(app.Listen(":3000"))
}
