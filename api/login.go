package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"seaottermsfs/model"
)

func Login(c *fiber.Ctx, store *session.Store, db *gorm.DB) error {
	var data model.LoginRequest
	var databaseData []model.User

	response := model.InitResponse()

	if err := c.BodyParser(&data); err != nil {
		logrus.Error(err)
		response.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	r := db.Model(&model.User{}).Find(&databaseData)
	if r.Error != nil {
		logrus.Error(r.Error)
		response.Message = r.Error.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	data.Username = strings.ToLower(data.Username)
	for _, col := range databaseData {
		if data.Username == col.Username {
			logrus.Infof("Username %s try to login", data.Username)
			err := bcrypt.CompareHashAndPassword([]byte(col.Password), []byte(data.Password))
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
				Username:   col.Username,
				Email:      col.Email,
				Avatar:     col.Avatar,
				Management: col.Management,
				CreatedAt:  col.CreatedAt,
				CreateName: col.CreateName,
			}

			return c.JSON(response)
		}
	}
	logrus.Error("user not found")
	response.Message = "user not found"
	return c.Status(fiber.StatusNotFound).JSON(response)
}

/* utils */
func CheckPassword(hashedPassword, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil
}
