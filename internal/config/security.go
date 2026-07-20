package config

import (
	"errors"
	"os"
	"runtime"
)

var (
	// ErrCredentialsNotFound is returned when the credentials file does not exist.
	ErrCredentialsNotFound = errors.New("credentials file does not exist")

	// ErrWeakerPermissions is returned when the credentials file permissions are weaker than 0600 on Unix systems.
	ErrWeakerPermissions = errors.New("security validation failed: credentials file permissions are too weak (must be strictly 0600, user read/write only)")
)

// VerifyCredentialsFile checks if the credentials file exists at the specified path.
// If it exists, it verifies that the file permissions on Unix systems are strictly 0600
// (Unix user read/write only, meaning group and other permissions are completely disabled, file mode & 0077 must be 0000).
// If the file does not exist, it returns ErrCredentialsNotFound.
func VerifyCredentialsFile(filePath string) error {
	fi, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrCredentialsNotFound
		}
		return err
	}

	if fi.IsDir() {
		return errors.New("credentials path is a directory, not a file")
	}

	// Permissions check on Unix systems
	if runtime.GOOS != "windows" {
		mode := fi.Mode().Perm()
		if mode&0077 != 0 {
			return ErrWeakerPermissions
		}
	}

	return nil
}
