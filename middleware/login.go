package middleware

import (
	"log/slog"
	"os"
	"seaottermsfs/model"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

// middleware for user login authentication
func LoginRequired(store *session.Store) fiber.Handler {
	return func(c fiber.Ctx) error {
		if !isLogin(c, store, false) {
			slog.Warn("user not login")
			return c.Status(fiber.StatusUnauthorized).JSON(model.GenerateResponse("請登入後再執行以下操作", nil))
		}
		return c.Next()
	}
}

// middleware for user login authentication(admin)
func LoginRequiredAdmin(store *session.Store) fiber.Handler {
	return func(c fiber.Ctx) error {
		if !isLogin(c, store, true) {
			slog.Warn("user is not admin")
			return c.Status(fiber.StatusUnauthorized).JSON(model.GenerateResponse("無效的操作", nil))
		}
		return c.Next()
	}
}

// check the user session
func isLogin(c fiber.Ctx, store *session.Store, isAdminRequired bool) bool {
	sess, err := store.Get(c)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	isAdmin, ok := sess.Get("isAdmin").(bool)
	if !ok || (isAdminRequired && !isAdmin) {
		return false
	}

	return true
}
