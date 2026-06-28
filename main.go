package main

import (
	"SimpleAV/config"
	"SimpleAV/database"
	"SimpleAV/engines"
	hashengine "SimpleAV/engines/hash"
	"SimpleAV/engines/quarantine"
	"SimpleAV/engines/watcher"
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"SimpleAV/applogger"
)

type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, ",") }
func (m *multiFlag) Set(val string) error {
	*m = append(*m, val)
	return nil
}

func main() {

	var wg sync.WaitGroup
	stop := make(chan os.Signal, 1)

	ctx, cancel := context.WithCancelCause(context.Background())
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// global configs
	config.Init()
	// logging configs
	applogger.Init(config.LogPath, false)
	// quarantine sir
	err := quarantine.Init()
	if err != nil {
		applogger.Error(err)
	}

	scan := flag.String("scan", "", "directory to scan, -scan path_to_dir")

	var watchPaths multiFlag
	flag.Var(&watchPaths, "watch", "directories to watch, -watch path_1 path_2...")

	flag.Parse()

	applogger.Info("AV Scanner up and running...")

	// TODO: Take args to read from csv and then set up DB
	err = database.ConnectDB()
	if err != nil {
		applogger.Error(err)
		os.Exit(1)
	}

	var sc = engines.NewScanner(hashengine.NewChecker(), watcher.NewWatcher(), quarantine.NewQuarantiner())
	// channle to feed scanner
	dirsToScan := make(chan string, 100)

	wg.Go(func() {
		for {
			select {
			case <-ctx.Done():
				return
			case file, ok := <-dirsToScan:
				if !ok {
					return
				}
				err := sc.ScanDirectory(file, ctx)
				if err != nil {
					applogger.Error(err)
				}
			}
		}
	})

	if *scan != "" {
		applogger.Info("scanning " + *scan)
		dirsToScan <- *scan
	}

	allPaths := append(config.DefaultWatchPaths(), watchPaths...)
	paths := getUniquePaths(&allPaths)

	if len(watchPaths) > 0 {
		applogger.Info("watching " + watchPaths.String())
	}
	wg.Go(func() {
		sc.Watch(paths, dirsToScan, ctx)
	})

	<-stop
	cancel(errors.New("SIGTERM"))
	close(dirsToScan)
	wg.Wait()
	applogger.Info("stopped, received stop signal")

	os.Exit(0)
}

func getUniquePaths(allPaths *[]string) []string {
	seen := make(map[string]struct{})
	deduped := make([]string, 0, len(*allPaths))

	for _, p := range *allPaths {
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		deduped = append(deduped, p)
	}
	return deduped
}
