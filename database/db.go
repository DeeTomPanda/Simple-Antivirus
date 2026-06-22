package database

import (
	"SimpleAV/apperrors"
	"SimpleAV/models"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() error {
	const dbName = "simpleAV.db"
	var err error
	var dbDir string

	switch runtime.GOOS {
	case "windows":
		dbDir = filepath.Join(os.Getenv("ProgramData"), "Simple-AV")
	case "linux":
		dbDir = filepath.Join(string(os.PathSeparator), "var", "lib", "simple_av")
	case "darwin":
		return fmt.Errorf("unsupported platform: %w", apperrors.ErrUnsupportedPlatform)
	}

	DB, err = OpenDB(dbDir, dbName)
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
	err := ensureDir(dbDir)
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

func ensureDir(dbDir string) error {
	return os.MkdirAll(dbDir, 0755)
}
