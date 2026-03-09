package router

import (
	"seaottermsfs/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func FolderRouter(routerGroup fiber.Router, store *session.Store) {
	fileGroup := routerGroup.Group("/file")

	fileGroup.Post("/*", func(c fiber.Ctx) error {
		path := c.Params("*")
		return service.CreateFolder(c, path)
	})
}
