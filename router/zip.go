package router

import (
	"seaottermsfs/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func ZipRouter(routerGroup fiber.Router, store *session.Store) {
	zipGroup := routerGroup.Group("/zip")

	zipGroup.Post("/", func(c fiber.Ctx) error {
		return service.ZipAllFiles(c)
	})

	zipGroup.Post("/:folderName", func(c fiber.Ctx) error {
		return service.ZipFiles(c)
	})
}
