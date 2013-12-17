package osutils

import (
	"os"
)

func PathExists(fileOrDirPath string) bool {
	if _, err := os.Stat(fileOrDirPath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
