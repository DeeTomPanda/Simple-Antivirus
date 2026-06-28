package sysutils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func EnsureDir(dbDir string) error {
	return os.MkdirAll(dbDir, 0755)
}

func ConvertToSHA256FromBytes(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

func ConvertToSHA256FromFile(file *os.File) (string, error) {

	hasher := sha256.New()

	// always reset to start of file
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	// get a string from raw byte output
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
