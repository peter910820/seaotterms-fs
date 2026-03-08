package router

import (
	"seaottermsfs/middleware"
	"seaottermsfs/service"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func ZipRouter(routerGroup fiber.Router, store *session.Store) {
	zipGroup := routerGroup.Group("/zip")

	zipGroup.Post("/", middleware.LoginRequired(store), func(c fiber.Ctx) error {
		return service.ZipAllFiles(c)
	})

	zipGroup.Post("/:folderName", middleware.LoginRequired(store), func(c fiber.Ctx) error {
		folderName := strings.TrimSpace(c.Params("folderName"))
		return service.ZipFiles(c, folderName)
	})
}
