package api

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"seaotterms-db/auth"
	"seaottermsfs/config"
	"seaottermsfs/model"
)

func Login(c *fiber.Ctx, store *session.Store) error {
	var data model.LoginRequest

	response := model.InitResponse()

	if err := c.BodyParser(&data); err != nil {
		logrus.Error(err)
		response.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	userData, err := auth.FindUserByUsername(config.Dbs.DB, strings.ToLower(data.Username))
	if err != nil {
		// user record not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Error("user not found")
			response.Message = "user not found"
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			logrus.Error(err)
			response.Message = err.Error()
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	logrus.Infof("Username %s try to login", data.Username)
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(data.Password))
	if err != nil {
		logrus.Error("login error, password not correct")
		response.Message = err.Error()
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}
	// set session
	sess, err := store.Get(c)
	if err != nil {
		logrus.Fatal(err)
	}
	sess.Set("username", data.Username)
	if err := sess.Save(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("Username %s login success", data.Username)

	response.Message = "Login Success"
	response.Data = model.LoginResponse{
		Username:   userData.Username,
		Email:      userData.Email,
		Avatar:     userData.Avatar,
		Management: userData.IsAdmin,
		CreatedAt:  userData.CreatedAt,
	}

	return c.JSON(response)
}

/* utils */
func CheckPassword(hashedPassword, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil
}
