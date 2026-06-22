package applogger

import (
	"SimpleAV/apperrors"
	"fmt"
	"log"
	"os"
)

var file *os.File
var logg *log.Logger
var debugMode bool

func Init(path string, debug bool) error {
	debugMode = debug
	var err error
	file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error on opening file")
	}

	logg = log.New(file, "AV-ENGINE: ", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}

func logger(level string, msg string, err error) {
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
	logger("INFO", msg, nil)
}
func Warn(msg string) {
	logger("WARN", msg, nil)
}
func Error(err error) {
	logger("ERROR", "", err)
}
func Debug() {}
