package main

import (
	"SimpleAV/config"
	"SimpleAV/database"
	"SimpleAV/engines"
	hashengine "SimpleAV/engines/hash"
	"SimpleAV/engines/watcher"
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"SimpleAV/applogger"
)

func main() {

	var wg sync.WaitGroup
	stop := make(chan os.Signal, 1)

	ctx, cancel := context.WithCancelCause(context.Background())
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	scan := flag.String("scan", "", "directory to scan, -scan path_to_dir")
	watch := flag.String("watch", "", "directory to watch -watch path_to_dir")

	flag.Parse()

	// global configs
	config.Init()
	// logging configs
	applogger.Init(config.LogPath, !false)

	applogger.Info("AV Scanner up and running...")

	// TODO: Take args to read from csv and then set up DB
	err := database.ConnectDB()
	if err != nil {
		applogger.Error(err)
		os.Exit(1)
	}

	if *scan != "" {
		applogger.Info("scanning " + *scan)
		wg.Go(func() {
			var sc = engines.NewScanner(&hashengine.Checker{})
			err = sc.ScanDirectory(*scan, ctx)
			if err != nil {
				applogger.Error(err)
			}
		})
	}

	if *watch != "" {
		applogger.Info("watching " + *watch)
		wg.Go(func() {
			watcher.Watch(ctx)
		})
	}

	<-stop
	cancel(errors.New("SIGTERM"))
	wg.Wait()
	applogger.Info("stopped, received stop signal")

	os.Exit(0)
}
