package service

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"

	"seaottermsfs/model"
	"seaottermsfs/utils"
)

const MaxUploadSize = 500 * 1024 * 1024 // 500MB

func Upload(c fiber.Ctx) error {
	directory, err := url.PathUnescape(strings.TrimSpace(c.FormValue("directory")))
	if err != nil {
		slog.Warn("Upload API失敗: 無效的路徑格式編碼")
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的路徑格式", nil))
	}
	directory = strings.TrimSpace(directory)
	if !utils.IsValidPathStructure(directory) {
		slog.Warn("Upload API失敗: 無效的路徑格式")
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的路徑格式", nil))
	}

	file, err := c.FormFile("file")
	if err != nil {
		slog.Warn("Upload API失敗: 讀取上傳檔案失敗: " + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("上傳失敗", nil))
	}

	rootPath := os.Getenv("RESOURCE_PATH")
	if rootPath == "" {
		slog.Error("Upload API失敗: 伺服器設定錯誤, RESOURCE_PATH is not set")
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("上傳失敗", nil))
	}
	rootPath = filepath.Clean(rootPath)

	directory = filepath.Clean(filepath.ToSlash(directory))
	if directory == "." {
		directory = ""
	}

	// prevent path traversal
	isPathSafe, targetDir := utils.IsPathSafe(rootPath, directory)
	if !isPathSafe {
		slog.Warn("Upload API失敗: 路徑安全檢查未通過")
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的檔案路徑", nil))
	}

	// filename
	baseName := strings.TrimSpace(c.FormValue("filename"))
	if baseName == "" {
		baseName = filepath.Base(file.Filename)
	} else {
		baseName, err = url.PathUnescape(baseName)
		if err != nil {
			slog.Warn("Upload API失敗: 無效的檔名編碼")
			return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的檔名", nil))
		}
		baseName = filepath.Base(baseName)
	}
	if baseName == "" || baseName == "." {
		slog.Warn("Upload API失敗: 無效的檔名")
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的檔名", nil))
	}
	if strings.Contains(baseName, "..") {
		slog.Warn("Upload API失敗: 檔名包含非法字元")
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的檔名", nil))
	}

	// prevent path traversal
	// verify that the parsed path is actually inside rootPath
	targetPath := filepath.Join(targetDir, baseName)
	relPath, err := filepath.Rel(rootPath, targetPath)
	if err != nil || strings.Contains(relPath, "..") {
		slog.Warn("Upload API失敗: 無效的儲存路徑")
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的儲存路徑", nil))
	}

	// file size check
	if file.Size > MaxUploadSize {
		slog.Warn("Upload API失敗: 檔案超過規定上限")
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("檔案超過規定上限", nil))
	}

	// directory must already exist and be a directory; do not create folders
	info, err := os.Stat(targetDir)
	if err != nil || !info.IsDir() {
		if err != nil {
			slog.Error(err.Error())
		}
		slog.Warn("Upload API失敗: 目標路徑不存在或並非目錄")
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("目標路徑不存在或並非目錄", nil))
	}

	// do not overwrite if a file with the same name already exists at that path
	if _, err := os.Stat(targetPath); err == nil {
		slog.Warn("Upload API失敗: 該路徑下已存在同名檔案")
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("該路徑下已存在同名檔案", nil))
	}

	if err := c.SaveFile(file, targetPath); err != nil {
		slog.Error("Upload API失敗: 儲存檔案失敗: " + err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse(err.Error(), nil))
	}

	slog.Info(fmt.Sprintf("%s 上傳成功，大小為 %d Bytes", baseName, file.Size))
	return c.Status(fiber.StatusOK).JSON(model.GenerateResponse("上傳成功", nil))
}
