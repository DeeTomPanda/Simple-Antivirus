package hash

import (
	"SimpleAV/apperrors"
	"SimpleAV/database"
	"SimpleAV/models"
	sysutils "SimpleAV/sys_utils"
	"errors"
	"fmt"
	"os"

	"gorm.io/gorm"
)

type Checker struct{}

func NewChecker() *Checker {
	return &Checker{}
}

func (c *Checker) CheckMaliciousHash(path string) (bool, error) {

	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// lock file to avoid modifications during scan
	err = sysutils.LockFile(file)
	if err != nil {
		return false, fmt.Errorf("file lock failed, %w:%w", apperrors.ErrLocking, err)
	}
	defer sysutils.UnlockFile(file)

	sha256Hash, err := sysutils.ConvertToSHA256FromFile(file)
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
