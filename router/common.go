package router

import (
	"seaottermsfs/api"

	"github.com/gofiber/fiber/v2"
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
