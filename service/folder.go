package service

import (
	"seaottermsfs/model"

	"github.com/gofiber/fiber/v3"
)

func CreateFolder(c fiber.Ctx, subPath string) error {
	return c.Status(fiber.StatusOK).JSON(model.GenerateResponse("success", nil))
}
