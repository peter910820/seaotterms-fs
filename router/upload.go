package router

import (
	"seaottermsfs/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func UploadRouter(routerGroup fiber.Router, store *session.Store) {
	uploadGroup := routerGroup.Group("/upload")

	uploadGroup.Post("/", func(c fiber.Ctx) error {
		return service.UploadFileV2(c)
	})
}
