package router

import (
	"seaottermsfs/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func FileRouter(routerGroup fiber.Router, store *session.Store) {
	fileGroup := routerGroup.Group("/file")

	fileGroup.Get("/*", func(c fiber.Ctx) error {
		path := c.Params("*")
		return service.GetFilesV2(c, path)
	})
}
