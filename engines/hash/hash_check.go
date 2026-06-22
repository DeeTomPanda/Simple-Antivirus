package hash

import (
	"SimpleAV/apperrors"
	"SimpleAV/database"
	"SimpleAV/models"
	sysutils "SimpleAV/sys_utils"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"gorm.io/gorm"
)

type Checker struct{}

func (c *Checker) CheckMaliciousHash(path string) (bool, error) {
	return CheckMaliciousHash(path)
}

func CheckMaliciousHash(filePath string) (bool, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// lock file to avoid modifications during scan
	err = sysutils.LockFile(file)
	defer sysutils.UnlockFile(file)

	sha256Hash, err := convertToSHA256(file)
	if err != nil {
		return false, fmt.Errorf("hashing err %w:%w", apperrors.ErrHashing, err)
	}

	exists, err := checkHashInDB(sha256Hash)
	if err != nil {
		return false, err
	}

	return exists, nil

}

func checkHashInDB(hash string) (bool, error) {

	var malware models.Malware
	result := database.DB.Where("sha256 = ?", hash).First(&malware)

	if result.Error != nil {
		// record not found, potentially not a malware
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("something wrong with DB %w : %w", apperrors.ErrDatabaseDown, result.Error)
	}

	// is malware !
	return true, nil

}

func convertToSHA256(file *os.File) (string, error) {

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
