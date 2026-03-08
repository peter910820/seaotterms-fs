package service

import (
	"archive/zip"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"seaottermsfs/model"
	"seaottermsfs/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

func ZipFiles(c fiber.Ctx) error {
	var data model.ZipRequest

	if err := c.Bind().Body(&data); err != nil {
		slog.Error(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse(err.Error(), nil))
	}

	folderName := strings.TrimSpace(data.FloderName)
	sourcePath := filepath.Join(os.Getenv("FILE_PATH"), folderName)

	// check target folder is exist
	folderExists, err := utils.CheckFolderExists(sourcePath)
	if !folderExists {
		if err != nil {
			slog.Error(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse("後端系統異常，請聯絡管理員", nil))
		}
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("壓縮的目標資料夾不存在", nil))
	}

	// ensure zip folder is exist
	zipDir := filepath.Join(os.Getenv("FILE_PATH"), "zip")
	if err = os.MkdirAll(zipDir, 0755); err != nil {
		slog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse("後端異常，請聯繫管理員", nil))
	}

	// build zip filename: folderName-yyyyMMddHHmm.zip
	timestamp := time.Now().Format("200601021504")
	zipFileName := folderName + "-" + timestamp + ".zip"
	zipFilePath := filepath.Join(zipDir, zipFileName)

	newZipFile, err := os.Create(zipFilePath)
	if err != nil {
		slog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse("後端異常，請聯繫管理員", nil))
	}
	defer newZipFile.Close()

	archive := zip.NewWriter(newZipFile)
	defer archive.Close()

	// walk through all files in source folder and add them to zip
	err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip directories, only add files
		if info.IsDir() {
			return nil
		}

		// build relative path inside the zip archive
		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}

		// read all file metadata (permissions, timestamps, etc.) from the OS
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// override specific fields as needed
		header.Name = filepath.ToSlash(relPath)
		header.Method = zip.Deflate // enable compression

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
	if err != nil {
		slog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse("壓縮過程發生錯誤，請聯繫管理員", nil))
	}

	return c.Status(fiber.StatusOK).JSON(model.GenerateResponse("壓縮成功", zipFileName))
}
