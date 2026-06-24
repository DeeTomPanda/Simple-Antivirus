package watcher

import (
	"SimpleAV/applogger"
	"context"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct{}

func (w *Watcher) Watch(paths []string, input chan<- string, ctx context.Context) {
	Watch(paths, input, ctx)
}

func Watch(paths []string, input chan<- string, ctx context.Context) error {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer fw.Close()

	// add all paths
	for _, path := range paths {
		if err := fw.Add(path); err != nil {
			applogger.Warn("could not watch: " + path)
		}
	}

	applogger.Info("watching paths...")

	for {
		select {
		case <-ctx.Done():
			applogger.Info("watcher stopping...")
			return nil

		case event, ok := <-fw.Events:
			if !ok {
				return nil
			}
			// only notify new/modified files
			if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) {
				applogger.Info("file event: " + event.Name)
				// run scan against it..
				select {
				case input <- event.Name:
				case <-ctx.Done():
					return nil
				}

			}
			if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				applogger.Warn("file removed/renamed: " + event.Name)
			}

		case err, ok := <-fw.Errors:
			if !ok {
				return nil
			}
			applogger.Error(err)
		}
	}
}
