package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"

	"seaottermsfs/config"
	"seaottermsfs/router"

	seaottermsdb "seaotterms-db"
)

var (
	// set frontendFolder
	frontendFolder string = "./dist"
	// init store(session)
	store = session.New(session.Config{
		Expiration: 12 * time.Hour,
		// CookieHTTPOnly: true,
	})
)

func init() {
	os.MkdirAll("./resource", os.ModePerm)
	os.MkdirAll("./resource/image", os.ModePerm)
	os.MkdirAll("./resource/test", os.ModePerm)
	err := godotenv.Load()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	w := os.Stderr
	logger := slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Stamp,
		}),
	)
	slog.SetDefault(logger)
}

func main() {
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// init db connection
	config.Dbs, err = seaottermsdb.InitDsn(seaottermsdb.ConnectDBConfig{
		Owner:    os.Getenv("DB_OWNER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     dbPort,
	})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// migration
	seaottermsdb.Migration(config.Dbs)

	app := fiber.New(
		fiber.Config{
			BodyLimit: 20 * 1024 * 1024, // 20MB
		})

	app.Static("/resource", "./resource")

	app.Use(cors.New(cors.Config{AllowOrigins: "http://localhost:8080",
		AllowMethods: "POST"}))

	// static folder
	app.Static("/", frontendFolder)

	// route group
	apiGroup := app.Group("/api") // main api route group
	// api router
	router.ApiRouter(apiGroup, store)
	router.LoginRouter(apiGroup, store)

	/* --------------------------------- */
	// match all routes
	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile(frontendFolder + "/index.html")
	})

	addr := fmt.Sprintf("127.0.0.1:%s", os.Getenv("PORT"))
	if err := app.Listen(addr); err != nil {
		slog.Error("start fiber server error", "address", addr, "error", err)
		os.Exit(1)
	}
}
