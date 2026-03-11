package router

import (
	"seaottermsfs/middleware"
	"seaottermsfs/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func UploadRouter(routerGroup fiber.Router, store *session.Store) {
	uploadGroup := routerGroup.Group("/upload")

	uploadGroup.Post("/", middleware.LoginRequired(store), func(c fiber.Ctx) error {
		return service.Upload(c)
	})
}
