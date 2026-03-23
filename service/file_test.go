package service

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"seaottermsfs/model"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func setupTestApp() *fiber.App {
	app := fiber.New()
	app.Delete("/file/*", func(c fiber.Ctx) error {
		path := c.Params("*")
		return DeleteFile(c, path)
	})
	return app
}

func TestDeleteFile(t *testing.T) {
	// Create a temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "seaottermsfs-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set the RESOURCE_PATH environment variable
	os.Setenv("RESOURCE_PATH", tmpDir)
	defer os.Unsetenv("RESOURCE_PATH")

	// Set up some test files and directories
	validFile := filepath.Join(tmpDir, "test-file.txt")
	if err := os.WriteFile(validFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	testDir := filepath.Join(tmpDir, "test-dir")
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatalf("failed to create test dir: %v", err)
	}

	app := setupTestApp()

	tests := []struct {
		name         string
		path         string // Path to delete
		expectedCode int
		setupFunc    func() // optional setup
	}{
		{
			name:         "Success deleting valid file",
			path:         "/file/test-file.txt",
			expectedCode: 200, // StatusOK
		},
		{
			name:         "Error file not found",
			path:         "/file/non-existent.txt",
			expectedCode: 404, // StatusNotFound
		},
		{
			name:         "Error trying to delete directory",
			path:         "/file/test-dir",
			expectedCode: 400, // StatusBadRequest
		},
		{
			name:         "Success deleting percent-encoded file path",
			path:         "/file/folder%20A/test%20space.txt",
			expectedCode: 200,
			setupFunc: func() {
				targetDir := filepath.Join(tmpDir, "folder A")
				if err := os.MkdirAll(targetDir, 0755); err != nil {
					t.Fatalf("failed to create encoded test dir: %v", err)
				}
				targetFile := filepath.Join(targetDir, "test space.txt")
				if err := os.WriteFile(targetFile, []byte("encoded test"), 0644); err != nil {
					t.Fatalf("failed to create encoded test file: %v", err)
				}
			},
		},
		{
			name:         "Error path traversal",
			path:         "/file/../file_test.go",
			expectedCode: 400, // StatusBadRequest
		},
		{
			name:         "Error encoded path traversal",
			path:         "/file/%2e%2e/file_test.go",
			expectedCode: 400, // StatusBadRequest
		},
		{
			name:         "Error empty path",
			path:         "/file/",
			expectedCode: 400, // StatusBadRequest
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.setupFunc()
			}

			req := httptest.NewRequest("DELETE", tt.path, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, resp.StatusCode)
				
				// Optional: print response body for debugging
				bodyBytes, _ := io.ReadAll(resp.Body)
				var apiResp model.Response
				if len(bodyBytes) > 0 {
					json.Unmarshal(bodyBytes, &apiResp)
					t.Logf("Response body: %+v", apiResp)
				}
			}
		})
	}
}
