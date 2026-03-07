package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func GetDirectory(c *fiber.Ctx) error {
	dir := "./resource"
	var dirName []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Error(err)
			return err
		}
		if info.IsDir() {
			dirName = append(dirName, path)
		}
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": dirName,
	})
}

func GetFiles(c *fiber.Ctx) error {
	folder := c.Query("folder")
	dir := "./resource"
	fileName := []string{}

	if strings.Contains(folder, "..") || strings.Contains(folder, "/") {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "Invalid path",
		})
	}
	if folder != "" {
		dir = "./resource/" + strings.TrimSpace(folder)
	}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Error(err)
			return err
		}
		if !info.IsDir() {
			fileName = append(fileName, path)
		}
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fileName,
	})
}

func UploadFile(c *fiber.Ctx) error {
	directory := c.FormValue("directory")
	file, err := c.FormFile("file")
	if err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "uploaded failed",
		})
	}

	directoryCheck := strings.ReplaceAll(directory, "resource/", "")
	if strings.Contains(directoryCheck, "..") || strings.Contains(directoryCheck, "/") || strings.Contains(file.Filename, "..") || strings.Contains(file.Filename, "/") {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "Invalid path",
		})
	}
	// 20MB upper limit
	if file.Size > 20971520 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "uploaded failed",
		})
	}
	logrus.Info(fmt.Sprintf("%s 上傳成功，大小為 %d Bytes", file.Filename, file.Size))
	directory = strings.ReplaceAll(directory, "\\", "/")
	err = c.SaveFile(file, fmt.Sprintf("./%s/%s", directory, file.Filename))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": "上傳成功",
	})
}
