package sysutils

import "os"

func EnsureDir(dbDir string) error {
	return os.MkdirAll(dbDir, 0755)
}
