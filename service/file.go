package service

import (
	"os"
	"path/filepath"
	"seaottermsfs/model"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func GetFilesV2(c fiber.Ctx, subPath string) error {
	rootPath := os.Getenv("FILE_PATH")
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
