package router

import (
	"seaottermsfs/middleware"
	"seaottermsfs/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func FolderRouter(routerGroup fiber.Router, store *session.Store) {
	folderGroup := routerGroup.Group("/folder")

	folderGroup.Post("/*", middleware.LoginRequired(store), func(c fiber.Ctx) error {
		path := c.Params("*")
		return service.CreateFolder(c, path)
	})
}
