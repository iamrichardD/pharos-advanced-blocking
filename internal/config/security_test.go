package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestVerifyCredentialsFile(t *testing.T) {
	// 1. Test: File does not exist
	t.Run("File does not exist", func(t *testing.T) {
		nonExistentPath := filepath.Join(t.TempDir(), "non_existent_file.json")
		err := VerifyCredentialsFile(nonExistentPath)
		if !errors.Is(err, ErrCredentialsNotFound) {
			t.Fatalf("expected ErrCredentialsNotFound, got %v", err)
		}
	})

	// 2. Test: File exists with strict 0600 permissions
	t.Run("File exists with strict 0600 permissions", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "credentials.json")

		// Create file with 0600 permissions
		err := os.WriteFile(filePath, []byte(`{"key": "secret"}`), 0600)
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		// Verify permissions
		err = VerifyCredentialsFile(filePath)
		if err != nil {
			t.Fatalf("expected no error for 0600 permissions, got: %v", err)
		}
	})

	// 3. Test: File exists with weak permissions (e.g. 0644, 0755)
	t.Run("File exists with weak permissions", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("skipping Unix permission tests on Windows")
		}

		weakModes := []os.FileMode{0644, 0755, 0660, 0606, 0601}

		for _, mode := range weakModes {
			t.Run(mode.String(), func(t *testing.T) {
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "credentials.json")

				// Create file
				err := os.WriteFile(filePath, []byte(`{"key": "secret"}`), mode)
				if err != nil {
					t.Fatalf("failed to create temp file: %v", err)
				}

				// Explicitly chmod to ensure mode is set (sometimes umask affects os.WriteFile)
				err = os.Chmod(filePath, mode)
				if err != nil {
					t.Fatalf("failed to chmod temp file: %v", err)
				}

				err = VerifyCredentialsFile(filePath)
				if !errors.Is(err, ErrWeakerPermissions) {
					t.Fatalf("expected ErrWeakerPermissions for mode %04o, got: %v", mode, err)
				}
			})
		}
	})

	// 4. Test: Target is a directory
	t.Run("Target path is a directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		err := VerifyCredentialsFile(tmpDir)
		if err == nil {
			t.Fatal("expected error when path is a directory, got nil")
		}
	})
}
