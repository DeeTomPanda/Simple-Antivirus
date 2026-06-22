package applogger

import (
	"SimpleAV/apperrors"
	"SimpleAV/config"
	sysutils "SimpleAV/sys_utils"
	"fmt"
	"io"
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
	var (
		err   error
		flags int
	)
	flags = log.Ldate | log.Ltime
	if debug {
		flags |= log.Lshortfile
	}

	err = sysutils.EnsureDir(path)
	if err != nil {
		fmt.Println("error on file creation %w", err)
	}
	file, err = os.OpenFile(filepath.Join(path, config.LogFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error on opening file: %w", err)
	}

	logg = log.New(io.MultiWriter(file, os.Stdout), "AV-ENGINE: ", flags)
	return nil
}

func logger(level LogLevel, msg string, err error) {
	if logg == nil {
		return
	}

	line := ""
	if err != nil {

		appErr := apperrors.Map(err)
		switch debugMode {
		case true:
			line = fmt.Sprintf("[%s] code: %v | msg: %s | trace: %+v", level, appErr.Code, appErr.Message, appErr.Err)
		case false:
			line = fmt.Sprintf("[%s] %s:", level, appErr.Message)
		}
	} else {
		line = fmt.Sprintf("[%s] %s", level, msg)
	}
	// log a level deep, dont show logging line
	logg.Output(3, line)
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
