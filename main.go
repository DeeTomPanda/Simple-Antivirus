package main

import (
	"SimpleAV/config"
	"SimpleAV/database"
	"os"

	"SimpleAV/applogger"
)

func main() {
	// global configs
	config.Init()
	// logging configs
	applogger.Init(config.LogPath, false)
	// TODO: Take args to read from csv and then set up DB
	err := database.ConnectDB()
	if err != nil {
		applogger.Error(err)
		os.Exit(1)
	}
}
