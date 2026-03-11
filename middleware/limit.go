package middleware

import (
	"log/slog"
	"os"
	"seaottermsfs/model"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

var (
	userFrequency = map[string]int{}
)

func init() {
	userFrequency = make(map[string]int)
}

// 預計做API呼叫限制(還未實作完成)
func UserFrequencyLimit(store *session.Store) fiber.Handler {
	return func(c fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		username, ok := sess.Get("username").(string)
		// 沒登入跳過紀錄
		if !ok {
			return c.Next()
		}
		userFrequency[username]++
		if userFrequency[username] > 10 {
			return c.Status(fiber.StatusTooManyRequests).JSON(model.GenerateResponse("請稍後再試", nil))
		}
		return c.Next()
	}
}
