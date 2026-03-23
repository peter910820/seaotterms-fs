package service

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"seaotterms-db/auth"
	"seaottermsfs/config"
	"seaottermsfs/model"
)

func Login(c fiber.Ctx, store *session.Store) error {
	var data model.LoginRequest

	if err := c.Bind().Body(&data); err != nil {
		slog.Warn("Login API失敗: 請求格式錯誤: " + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse(err.Error(), nil))
	}

	userData, err := auth.FindUserByUsername(config.Dbs.DB, strings.ToLower(data.Username))
	if err != nil {
		// user record not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("Login API失敗: 使用者不存在")
			return c.Status(fiber.StatusNotFound).JSON(model.GenerateResponse("user not found", nil))
		} else {
			slog.Error("Login API失敗: 查詢使用者資料失敗: " + err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse(err.Error(), nil))
		}
	}

	slog.Info(fmt.Sprintf("Username %s try to login", data.Username))
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(data.Password))
	if err != nil {
		slog.Warn("Login API失敗: 密碼錯誤")
		return c.Status(fiber.StatusUnauthorized).JSON(model.GenerateResponse(err.Error(), nil))
	}
	// set session
	sess, err := store.Get(c)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	sess.Set("username", data.Username)
	sess.Set("isAdmin", userData.IsAdmin)
	if err := sess.Save(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info(fmt.Sprintf("Username %s login success", data.Username))

	loginResp := model.LoginResponse{
		Username:  userData.Username,
		Email:     userData.Email,
		Avatar:    userData.Avatar,
		IsAdmin:   userData.IsAdmin,
		CreatedAt: userData.CreatedAt,
	}

	return c.JSON(model.GenerateResponse("Login Success", loginResp))
}

/* utils */
func CheckPassword(hashedPassword, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil
}
