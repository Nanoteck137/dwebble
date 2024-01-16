package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/nanoteck137/dwebble/v2/internal/database"
	"github.com/nanoteck137/dwebble/v2/internal/handlers"
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

	app := fiber.New(fiber.Config{
		BodyLimit: 1 * 1024 * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if strings.HasPrefix(c.Path(), "/api") {
				switch err := err.(type) {
				case handlers.ApiError:
					return c.Status(err.Status).JSON(err)
				case *fiber.Error:
					return c.Status(err.Code).JSON(handlers.ApiError{
						Status:  err.Code,
						Message: err.Error(),
						Data:    nil,
					})
				default:
					return c.Status(500).JSON(handlers.ApiError{
						Status:  500,
						Message: err.Error(),
						Data:    nil,
					})
				}
			}

			return fiber.DefaultErrorHandler(c, err)
		},
	})

	app.Use(logger.New())

	app.Use(cors.New())

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

	apiConfig := handlers.New(db, queries, "./work")

	v1 := app.Group("/api/v1")

	v1.Get("/artists", apiConfig.HandlerGetAllArtists)
	v1.Post("/artists", apiConfig.HandlerCreateArtist)
	v1.Get("/artists/:id", apiConfig.HandlerGetArtist)

	v1.Get("/albums", apiConfig.HandlerGetAllAlbums)
	v1.Post("/albums", apiConfig.HandlerCreateAlbum)
	v1.Get("/albums/:id", apiConfig.HandlerGetAlbum)

	v1.Get("/tracks", apiConfig.HandlerGetAllTracks)
	v1.Post("/tracks", apiConfig.HandlerCreateTrack)
	v1.Get("/tracks/:id", apiConfig.HandlerGetTrack)

	// /queue/album/:id
	v1.Post("/queue/album", apiConfig.HandlerCreateQueueFromAlbum)

	// v1.Post("/", func(c *fiber.Ctx) error {
	// 	form, err  := c.MultipartForm()
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	t := form.File["files"]
	// 	for _, x := range t {
	// 		fmt.Printf("%v\n", x.Filename)
	// 	}
	//
	// 	return nil
	// })

	log.Fatal(app.Listen(":3000"))
}
