package config

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	LogPath string
	DBPath  string
)

func Init() {
	switch runtime.GOOS {
	case "windows":
		base := os.Getenv("ProgramData")
		LogPath = filepath.Join(base, "Simple-AV", "AV.log")
		DBPath = filepath.Join(base, "Simple-AV", "simpleAV.db")
	case "linux":
		LogPath = filepath.Join("/", "var", "log", "simple_av", "AV.log")
		DBPath = filepath.Join("/", "var", "lib", "simple_av", "simpleAV.db")
	}
}
