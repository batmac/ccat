package utils

import "os"

func PathExists(path string) bool {
	_, err := os.Stat(ExpandPath(path))
	return err == nil
}
