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

// zip target floder files
func ZipFiles(c fiber.Ctx, folderName string) error {
	if folderName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("請提供資料夾名稱", nil))
	}
	rootPath := os.Getenv("RESOURCE_PATH")
	folderName = filepath.Clean(folderName)

	// prevent path traversal
	// verify that the parsed path is actually inside rootPath
	sourcePath := filepath.Join(rootPath, folderName)
	rel, err := filepath.Rel(rootPath, sourcePath)
	if err != nil || strings.Contains(rel, "..") || rel == "." {
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("無效的資料夾路徑", nil))
	}

	// check target folder is exist
	folderExists, err := utils.CheckFolderExists(sourcePath)
	if !folderExists {
		if err != nil {
			slog.Error(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse("後端系統異常，請聯絡管理員", nil))
		}
		return c.Status(fiber.StatusBadRequest).JSON(model.GenerateResponse("壓縮的目標資料夾不存在", nil))
	}

	timestamp := time.Now().Format("200601021504")
	baseName := filepath.Base(sourcePath)
	zipFileName := baseName + "-" + timestamp + ".zip"

	return handleZipCreation(c, rootPath, sourcePath, zipFileName)
}

// zip all floder files
func ZipAllFiles(c fiber.Ctx) error {
	rootPath := os.Getenv("RESOURCE_PATH")

	timestamp := time.Now().Format("200601021504")
	zipFileName := "root-" + timestamp + ".zip"

	return handleZipCreation(c, rootPath, rootPath, zipFileName)
}

// handles the actual zip creation process.
// rootDir is the base directory of the system, and targetDir is the directory to compress.
func handleZipCreation(c fiber.Ctx, rootDir, targetDir, zipFileName string) error {
	zipDir := filepath.Join(rootDir, "zip")
	if err := os.MkdirAll(zipDir, 0755); err != nil {
		slog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse("後端異常，請聯繫管理員", nil))
	}

	zipFilePath := filepath.Join(zipDir, zipFileName)

	newZipFile, err := os.Create(zipFilePath)
	if err != nil {
		slog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.GenerateResponse("後端異常，請聯繫管理員", nil))
	}
	defer newZipFile.Close()

	archive := zip.NewWriter(newZipFile)
	defer archive.Close()

	err = filepath.WalkDir(targetDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// skip the target folder itself so we don't create an empty top-level folder entry
		if path == targetDir {
			return nil
		}

		// skip the zip directory to avoid recursive zipping when targetDir == rootDir
		if d.IsDir() && path == zipDir {
			return filepath.SkipDir
		}

		// skip symlinks
		if d.Type()&os.ModeSymlink != 0 {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// build relative path inside the zip archive
		relPath, err := filepath.Rel(targetDir, path)
		if err != nil {
			return err
		}

		info, err := d.Info()
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
		if d.IsDir() {
			// ensure trailing slash for directories
			if !strings.HasSuffix(header.Name, "/") {
				header.Name += "/"
			}
			header.Method = zip.Store // no compression for directories
			_, err = archive.CreateHeader(header)
			return err
		}

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
