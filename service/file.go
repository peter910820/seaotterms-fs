package service

import (
	"log/slog"
	"os"
	"path/filepath"
	"seaottermsfs/model"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func GetFiles(c fiber.Ctx, subPath string) error {
	rootPath := os.Getenv("RESOURCE_PATH")
	files := []string{}
	directories := []string{}

	// prevent path traversal
	// verify that the parsed path is actually inside rootPath
	sourcePath := filepath.Join(rootPath, filepath.Clean(subPath))
	rel, err := filepath.Rel(rootPath, sourcePath)
	if err != nil || strings.Contains(rel, "..") {
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的資料夾路徑", nil))
	}

	// traverse the current level only
	entries, err := os.ReadDir(sourcePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse("無法讀取目錄: "+err.Error(), nil))
	}

	for _, entry := range entries {
		// skip hidden files
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		if entry.IsDir() {
			directories = append(directories, entry.Name())
		} else {
			files = append(files, entry.Name())
		}
	}

	result := model.FileResponse{
		Files:       files,
		Directories: directories,
	}

	return c.Status(fiber.StatusOK).JSON(model.GenerateResponse("success", result))
}

// delete a file at the given path (relative to RESOURCE_PATH)
//
// Path must be a file, not a directory.
func DeleteFile(c fiber.Ctx, subPath string) error {
	rootPath := os.Getenv("RESOURCE_PATH")
	if rootPath == "" {
		slog.Error("RESOURCE_PATH is not set")
		return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse("伺服器設定錯誤", nil))
	}
	rootPath = filepath.Clean(rootPath)

	subPath = filepath.Clean(strings.ReplaceAll(strings.TrimSpace(subPath), "\\", "/"))
	if subPath == "" || subPath == "." || strings.Contains(subPath, "..") {
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的檔案路徑", nil))
	}

	targetPath := filepath.Join(rootPath, subPath)
	rel, err := filepath.Rel(rootPath, targetPath)
	if err != nil || strings.Contains(rel, "..") {
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的檔案路徑", nil))
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return c.Status(fiber.StatusNotFound).JSON(model.GenerateResponse("檔案不存在", nil))
		}
		slog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse("無法讀取檔案", nil))
	}
	if info.IsDir() {
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("此路徑為資料夾，僅可刪除檔案", nil))
	}

	if err := os.Remove(targetPath); err != nil {
		slog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse("刪除失敗", nil))
	}
	slog.Info("檔案已刪除: " + subPath)
	return c.Status(fiber.StatusOK).JSON(model.GenerateResponse("已刪除", nil))
}
