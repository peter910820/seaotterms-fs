package service

import (
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"seaottermsfs/model"
	"seaottermsfs/utils"

	"github.com/gofiber/fiber/v3"
)

// 建立資料夾
//
// 若父層不存在會一併建立
func CreateFolder(c fiber.Ctx, subPath string) error {
	rootPath := os.Getenv("RESOURCE_PATH")
	if rootPath == "" {
		slog.Error("CreateFolder API失敗: 伺服器設定錯誤, RESOURCE_PATH is not set")
		return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse("伺服器設定錯誤", nil))
	}
	rootPath = filepath.Clean(rootPath)

	decodedPath, err := url.PathUnescape(strings.TrimSpace(subPath))
	if err != nil {
		slog.Warn("CreateFolder API失敗: 無效的資料夾路徑編碼")
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的資料夾路徑", nil))
	}

	subPath = filepath.Clean(filepath.ToSlash(decodedPath))
	if subPath == "" || subPath == "." || strings.Contains(subPath, "..") {
		slog.Warn("CreateFolder API失敗: 無效的資料夾路徑")
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的資料夾路徑", nil))
	}

	isPathSafe, targetPath := utils.IsPathSafe(rootPath, subPath)
	if !isPathSafe {
		slog.Warn("CreateFolder API失敗: 路徑安全檢查未通過")
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的資料夾路徑", nil))
	}

	info, err := os.Stat(targetPath)
	if err == nil {
		if info.IsDir() {
			slog.Info("CreateFolder API成功: 資料夾已存在 " + subPath)
			return c.Status(fiber.StatusOK).JSON(model.GenerateResponse("資料夾已存在", nil))
		}
		slog.Warn("CreateFolder API失敗: 路徑已存在且為檔案 " + subPath)
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("該路徑已存在且為檔案，無法建立資料夾", nil))
	}
	if !os.IsNotExist(err) {
		slog.Error("CreateFolder API失敗: 無法讀取路徑: " + err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse("無法讀取路徑", nil))
	}

	if err := os.MkdirAll(targetPath, 0750); err != nil {
		slog.Error("CreateFolder API失敗: 建立資料夾失敗: " + err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse("建立資料夾失敗", nil))
	}
	slog.Info("資料夾已建立: " + subPath)
	return c.Status(fiber.StatusOK).JSON(model.GenerateResponse("資料夾已建立", nil))
}
