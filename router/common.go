package router

import (
	"seaottermsfs/middleware"
	"seaottermsfs/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func ApiRouter(routerGroup fiber.Router, store *session.Store) {
	routerGroup.Get("/directory", middleware.LoginRequired(store), func(c fiber.Ctx) error {
		return service.GetDirectory(c)
	})
	routerGroup.Get("/files", middleware.LoginRequired(store), func(c fiber.Ctx) error {
		return service.GetFiles(c)
	})
	routerGroup.Post("/upload", middleware.LoginRequired(store), func(c fiber.Ctx) error {
		return service.UploadFile(c)
	})
}
