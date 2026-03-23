package service

import (
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"seaottermsfs/model"
	"seaottermsfs/utils"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func GetFiles(c fiber.Ctx, subPath string) error {
	rootPath := os.Getenv("RESOURCE_PATH")
	files := []string{}
	directories := []string{}

	decodedPath, err := url.PathUnescape(strings.TrimSpace(subPath))
	if err != nil {
		slog.Warn("GetFiles API失敗: 無效的檔案路徑, rawPath=" + subPath)
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的檔案路徑", nil))
	}

	// prevent path traversal
	isPathSafe, sourcePath := utils.IsPathSafe(rootPath, filepath.Clean(filepath.ToSlash(decodedPath)))
	if !isPathSafe {
		slog.Warn("GetFiles API失敗: 路徑安全檢查未通過, rawPath=" + subPath)
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的檔案路徑", nil))
	}

	// traverse the current level only
	entries, err := os.ReadDir(sourcePath)
	if err != nil {
		slog.Error("GetFiles API失敗: 無法讀取目錄: " + err.Error())
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

	if strings.TrimSpace(decodedPath) == "" {
		slog.Info("GetFiles API成功: ./")
	} else {
		slog.Info("GetFiles API成功: " + decodedPath)
	}
	return c.Status(fiber.StatusOK).JSON(model.GenerateResponse("success", result))
}

// delete a file at the given path (relative to RESOURCE_PATH)
//
// Path must be a file, not a directory.
func DeleteFile(c fiber.Ctx, subPath string) error {
	rootPath := os.Getenv("RESOURCE_PATH")
	if rootPath == "" {
		slog.Error("DeleteFile API失敗: 伺服器設定錯誤, RESOURCE_PATH is not set")
		return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse("伺服器設定錯誤", nil))
	}
	rootPath = filepath.Clean(rootPath)

	decodedPath, err := url.PathUnescape(strings.TrimSpace(subPath))
	if err != nil {
		slog.Warn("DeleteFile API失敗: 無效的檔案路徑編碼")
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的檔案路徑", nil))
	}

	subPath = filepath.Clean(filepath.ToSlash(decodedPath))
	if subPath == "" || subPath == "." || strings.Contains(subPath, "..") {
		slog.Warn("DeleteFile API失敗: 無效的檔案路徑")
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的檔案路徑", nil))
	}

	// prevent path traversal
	isPathSafe, targetPath := utils.IsPathSafe(rootPath, subPath)
	if !isPathSafe {
		slog.Warn("DeleteFile API失敗: 路徑安全檢查未通過")
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的檔案路徑", nil))
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn("DeleteFile API失敗: 檔案不存在 " + subPath)
			return c.Status(fiber.StatusNotFound).JSON(model.GenerateResponse("檔案不存在", nil))
		}
		slog.Error("DeleteFile API失敗: 無法讀取檔案: " + err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse("無法讀取檔案", nil))
	}
	if info.IsDir() {
		slog.Warn("DeleteFile API失敗: 目標為資料夾 " + subPath)
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("此路徑為資料夾，僅可刪除檔案", nil))
	}

	if err := os.Remove(targetPath); err != nil {
		slog.Error("DeleteFile API失敗: 刪除失敗: " + err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse("刪除失敗", nil))
	}
	slog.Info("檔案已刪除: " + subPath)
	return c.Status(fiber.StatusOK).JSON(model.GenerateResponse("已刪除", nil))
}
