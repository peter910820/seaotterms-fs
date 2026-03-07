package router

import (
	"os"

	"seaottermsfs/api"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func ApiRouter(routerGroup fiber.Router) {
	routerGroup.Get("/directory", func(c *fiber.Ctx) error {
		return api.GetDirectory(c)
	})
	routerGroup.Get("/files", func(c *fiber.Ctx) error {
		return api.GetFiles(c)
	})
	routerGroup.Post("/upload", func(c *fiber.Ctx) error {
		return api.UploadFile(c)
	})
}

func LoginRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	loginGroup := routerGroup.Group("/login")
	loginGroup.Post("/", func(c *fiber.Ctx) error {
		return api.Login(c, store, dbs[os.Getenv("DB_NAME")])
	})
}
