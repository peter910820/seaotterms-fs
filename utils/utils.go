package utils

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

func CheckFolderExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// prevent path traversal
//
// verify that the parsed path is actually inside rootPath
func IsPathSafe(rootPath, directory string) (save bool, targetDir string) {
	save = true
	targetDir = filepath.Join(rootPath, directory)
	relDir, err := filepath.Rel(rootPath, targetDir)
	if err != nil || strings.Contains(relDir, "..") {
		if err != nil {
			slog.Warn(err.Error())
		} else {
			slog.Warn("偵測到惡意路徑穿越嘗試: " + directory)
		}
		save = false
	}
	return
}

// 更嚴格的檢查(上傳API用)
func IsValidPathStructure(s string) bool {
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
