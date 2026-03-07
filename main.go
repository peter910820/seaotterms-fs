package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"seaottermsfs/config"
	"seaottermsfs/router"
)

var (
	// init store(session)
	store = session.New(session.Config{
		Expiration: 12 * time.Hour,
		// CookieHTTPOnly: true,
	})
	// management database connect
	dbs = make(map[string]*gorm.DB)
	// set frontendFolder
	frontendFolder string = "./dist"
)

func init() {
	os.MkdirAll("./resource", os.ModePerm)
	os.MkdirAll("./resource/image", os.ModePerm)
	os.MkdirAll("./resource/test", os.ModePerm)
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal(err)
	}
	// logrus settings
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	// init dsn
	dbName, db := config.InitDsn()
	dbs[dbName] = db

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
	router.ApiRouter(apiGroup)
	router.LoginRouter(apiGroup, store, dbs)

	/* --------------------------------- */
	// match all routes
	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile(frontendFolder + "/index.html")
	})

	logrus.Fatal(app.Listen(fmt.Sprintf("127.0.0.1:%s", os.Getenv("PORT"))))
}
