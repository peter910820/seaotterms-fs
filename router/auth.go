package router

import (
	"seaottermsfs/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func LoginRouter(routerGroup fiber.Router, store *session.Store) {
	loginGroup := routerGroup.Group("/login")
	loginGroup.Post("/", func(c *fiber.Ctx) error {
		return service.Login(c, store)
	})
}
