package watcher

import (
	"SimpleAV/apperrors"
	"SimpleAV/applogger"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const debounceInterval = 500 * time.Millisecond

type Watcher struct {
	mutex     sync.Mutex
	fileTimer map[string]*time.Timer
}

func NewWatcher() *Watcher {
	return &Watcher{
		fileTimer: make(map[string]*time.Timer),
	}
}

func (w *Watcher) debounce(filePath string, input chan<- string, ctx context.Context) {
	w.mutex.Lock()

	if fileTimer, ok := w.fileTimer[filePath]; ok {
		fileTimer.Stop()
	}

	// create a new timer and store it in the map
	w.fileTimer[filePath] = time.AfterFunc(debounceInterval, func() {
		select {
		case input <- filePath:
		case <-ctx.Done():
		}

		w.mutex.Lock()
		delete(w.fileTimer, filePath)
		w.mutex.Unlock()
	})

	w.mutex.Unlock()
}

func (w *Watcher) Watch(paths []string, input chan<- string, ctx context.Context) error {

	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("%w:%w", apperrors.ErrFileWatch, err)
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
				// run scan against it...
				w.debounce(event.Name, input, ctx)

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
