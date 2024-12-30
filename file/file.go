package file

import (
	"os"
)

// IsExistE returns whether this path exist and error
func IsExistE(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	return os.IsExist(err), err
}

// IsExist reports whether this path exist
func IsExist(path string) bool {
	b, _ := IsExistE(path)
	return b
}

// IsDirE returns an error as whether this path is a directory
func IsDirE(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return f.IsDir(), nil
}

// IsDir reports whether this path is a directory
func IsDir(path string) bool {
	b, _ := IsDirE(path)
	return b
}
