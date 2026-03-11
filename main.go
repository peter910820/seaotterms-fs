package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"

	"seaottermsfs/config"
	"seaottermsfs/model"
	"seaottermsfs/router"
	"seaottermsfs/service"

	seaottermsdb "seaotterms-db"
)

var (
	// init store(session)
	store = session.NewStore(session.Config{
		IdleTimeout: 6 * time.Hour,
		// CookieHTTPOnly: true,
	})
)

func init() {
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
			BodyLimit: service.MaxUploadSize,
		})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("ALLOW_ORIGINS")},
		AllowMethods:     []string{"POST", "GET", "DELETE"},
		AllowCredentials: true,
	}))

	// route group
	apiGroup := app.Group("/api") // main api route group
	// api router
	router.LoginRouter(apiGroup, store)
	router.FileRouter(apiGroup, store)
	router.FolderRouter(apiGroup, store)
	router.UploadRouter(apiGroup, store)
	router.ZipRouter(apiGroup, store)

	// match all routes
	app.Get("*", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(model.GenerateResponse("404 not found", nil))
	})

	// 每 12 小時清除過期 session (強制登出所有人)
	go func() {
		ticker := time.NewTicker(12 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			if err := store.Reset(ctx); err != nil {
				slog.Error("session cleanup failed", "error", err)
			} else {
				slog.Info("session store reset (expired sessions cleared)")
			}
			cancel()
		}
	}()

	addr := fmt.Sprintf("127.0.0.1:%s", os.Getenv("PORT"))
	if err := app.Listen(addr); err != nil {
		slog.Error("start fiber server error", "address", addr, "error", err)
		os.Exit(1)
	}
}
