package service

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/gofiber/fiber/v3"

	"seaottermsfs/model"
)

const MaxUploadSize = 500 * 1024 * 1024 // 500MB

func UploadFileV2(c fiber.Ctx) error {
	directory := strings.TrimSpace(c.FormValue("directory"))
	if !isValidPathStructure(directory) {
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的路徑格式", nil))
	}
	file, err := c.FormFile("file")
	if err != nil {
		slog.Error(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("上傳失敗", nil))
	}

	rootPath := os.Getenv("RESOURCE_PATH")
	if rootPath == "" {
		slog.Error("RESOURCE_PATH is not set")
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("上傳失敗", nil))
	}
	rootPath = filepath.Clean(rootPath)

	directory = strings.ReplaceAll(directory, "\\", "/")
	directory = filepath.Clean(directory)
	if directory == "." {
		directory = ""
	}

	// prevent path traversal
	// verify that the parsed path is actually inside rootPath
	targetDir := filepath.Join(rootPath, directory)
	relDir, err := filepath.Rel(rootPath, targetDir)
	if err != nil || strings.Contains(relDir, "..") {
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的資料夾路徑", nil))
	}

	// filename must not contain path information
	baseName := filepath.Base(file.Filename)
	if baseName == "" || baseName == "." {
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的檔名", nil))
	}
	if strings.Contains(baseName, "..") {
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的檔名", nil))
	}

	// prevent path traversal
	// verify that the parsed path is actually inside rootPath
	targetPath := filepath.Join(targetDir, baseName)
	relPath, err := filepath.Rel(rootPath, targetPath)
	if err != nil || strings.Contains(relPath, "..") {
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的儲存路徑", nil))
	}

	// file size check
	if file.Size > MaxUploadSize {
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("檔案超過規定上限", nil))
	}

	// directory must already exist and be a directory; do not create folders
	info, err := os.Stat(targetDir)
	if err != nil || !info.IsDir() {
		if err != nil {
			slog.Error(err.Error())
		}
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("目標路徑不存在或並非目錄", nil))
	}

	// do not overwrite if a file with the same name already exists at that path
	if _, err := os.Stat(targetPath); err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("該路徑下已存在同名檔案", nil))
	}

	if err := c.SaveFile(file, targetPath); err != nil {
		slog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse(err.Error(), nil))
	}

	slog.Info(fmt.Sprintf("%s 上傳成功，大小為 %d Bytes", baseName, file.Size))
	return c.Status(fiber.StatusOK).JSON(model.GenerateResponse("上傳成功", nil))
}

func isValidPathStructure(s string) bool {
	if !utf8.ValidString(s) {
		return false
	}
	if strings.Contains(s, "\x00") || strings.Contains(s, "..") {
		return false
	}

	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return false
	}
	for _, r := range trimmed {
		switch {
		case r == '/' || r == '-' || r == '_' || r == '.':
			continue
		case r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9':
			continue
		default:
			return false
		}
	}
	return true
}
