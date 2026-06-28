package config

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	DBName        = "simpleAV.db"
	LogFileName   = "AV.log"
	LogPath       string
	DBPath        string
	QuarantineDir string
)

func Init() {
	switch runtime.GOOS {
	case "windows":
		base := os.Getenv("ProgramData")
		LogPath = filepath.Join(base, "Simple-AV")
		DBPath = filepath.Join(base, "Simple-AV")
		QuarantineDir = filepath.Join(base, "Simple-AV", "quarantine")
	case "linux":
		LogPath = filepath.Join("/", "var", "log", "simple_av")
		DBPath = filepath.Join("/", "var", "lib", "simple_av")
		QuarantineDir = filepath.Join("/", "var", "lib", "simple_av", "quarantine")
	}
}

func DefaultWatchPaths() []string {
	switch runtime.GOOS {
	case "linux":
		return []string{
			"/home",
			"/etc",
			"/usr/bin",
			"/usr/local/bin",
			"/root",
		}
	case "windows":
		return []string{
			`C:\Windows\System32`,
			`C:\Windows\Temp`,
			`C:\Users`,
			`C:\ProgramData`,
		}
	default:
		return []string{}
	}
}
