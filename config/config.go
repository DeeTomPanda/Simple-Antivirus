package config

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	LogPath     string
	DBPath      string
	DBName      = "simpleAV.db"
	LogFileName = "AV.log"
)

func Init() {
	switch runtime.GOOS {
	case "windows":
		base := os.Getenv("ProgramData")
		LogPath = filepath.Join(base, "Simple-AV")
		DBPath = filepath.Join(base, "Simple-AV")
	case "linux":
		LogPath = filepath.Join("/", "var", "log", "simple_av")
		DBPath = filepath.Join("/", "var", "lib", "simple_av")
	}
}
