package router

import (
	"seaottermsfs/api"
	"seaottermsfs/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func ApiRouter(routerGroup fiber.Router, store *session.Store) {
	routerGroup.Get("/directory", middleware.LoginRequired(store), func(c *fiber.Ctx) error {
		return api.GetDirectory(c)
	})
	routerGroup.Get("/files", middleware.LoginRequired(store), func(c *fiber.Ctx) error {
		return api.GetFiles(c)
	})
	routerGroup.Post("/upload", middleware.LoginRequired(store), func(c *fiber.Ctx) error {
		return api.UploadFile(c)
	})
}
