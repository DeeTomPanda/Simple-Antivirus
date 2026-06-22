package database

import (
	"SimpleAV/apperrors"
	"SimpleAV/config"
	"SimpleAV/models"
	sysutils "SimpleAV/sys_utils"
	"fmt"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() error {
	var err error

	DB, err = OpenDB(config.DBPath, config.DBName)
	if err != nil {
		return err
	}

	// create schema
	err = DB.AutoMigrate(&models.Malware{})
	if err != nil {
		return fmt.Errorf("failed migration %w: %w", apperrors.ErrDatabaseDown, err)
	}
	return nil

}

func OpenDB(dbDir string, dbName string) (*gorm.DB, error) {
	err := sysutils.EnsureDir(dbDir)
	if err != nil {
		return nil, fmt.Errorf("ensure DBDir exists %w: %w", apperrors.ErrDatabaseDown, err)
	}

	db, err := gorm.Open(sqlite.Open(filepath.Join(dbDir, dbName)), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, fmt.Errorf("db failed to open %w: %w", apperrors.ErrDatabaseDown, err)
	}

	return db, err
}
