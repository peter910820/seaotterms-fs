package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

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
	// init db connection
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		logrus.Fatal(err)
	}
	config.Dbs, err = seaottermsdb.InitDsn(seaottermsdb.ConnectDBConfig{
		Owner:    os.Getenv("DB_OWNER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     dbPort,
	})
	if err != nil {
		logrus.Fatal(err)
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
	router.ApiRouter(apiGroup)
	router.LoginRouter(apiGroup, store)

	/* --------------------------------- */
	// match all routes
	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile(frontendFolder + "/index.html")
	})

	logrus.Fatal(app.Listen(fmt.Sprintf("127.0.0.1:%s", os.Getenv("PORT"))))
}
