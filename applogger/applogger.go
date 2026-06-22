package applogger

import (
	"SimpleAV/apperrors"
	"SimpleAV/config"
	sysutils "SimpleAV/sys_utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var file *os.File
var logg *log.Logger
var debugMode bool

type LogLevel string

var (
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

func Init(path string, debug bool) error {
	debugMode = debug
	var err error

	err = sysutils.EnsureDir(path)
	if err != nil {
		fmt.Println("error on file creation %w", err)
	}
	file, err = os.OpenFile(filepath.Join(path, config.LogFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error on opening file: %w", err)
	}

	logg = log.New(file, "AV-ENGINE: ", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}

func logger(level LogLevel, msg string, err error) {
	if logg == nil {
		return
	}

	if err != nil {
		appErr := apperrors.Map(err)
		switch debugMode {
		case true:
			logg.Printf("[%s] code: %v | msg: %s | trace: %+v", level, appErr.Code, appErr.Message, appErr.Err)
		case false:
			logg.Printf("[%s] %s:", level, appErr.Message)
		}

	} else {
		logg.Printf("[%s] %s", level, msg)
	}
}

func Info(msg string) {
	logger(INFO, msg, nil)
}
func Warn(msg string) {
	logger(WARN, msg, nil)
}
func Error(err error) {
	logger(ERROR, "", err)
}
func Debug() {}
