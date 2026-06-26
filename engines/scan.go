package engines

import (
	"SimpleAV/applogger"
	"context"
	"io/fs"
	"path/filepath"
)

type HashChecker interface {
	CheckMaliciousHash(path string) (bool, error)
}

type FileWatcher interface {
	Watch(paths []string, input chan<- string, ctx context.Context) error
}

type Scanner struct {
	hashChecker HashChecker
	fileWatcher FileWatcher
}

func NewScanner(hashChecker HashChecker, watcher FileWatcher) *Scanner {
	return &Scanner{
		hashChecker: hashChecker,
		fileWatcher: watcher,
	}
}

func (s *Scanner) ScanDirectory(root string, ctx context.Context) error {

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {

		if err := ctx.Err(); err != nil {
			return err
		}

		if err != nil {
			return err
		}

		// skip dir
		if d.IsDir() {
			return nil
		}

		malicious, err := s.hashChecker.CheckMaliciousHash(path)
		if err != nil {
			return err
		}

		if malicious {
			applogger.Warn("malware detected: " + path)
		} else {
			applogger.Info("Clean! " + path)
		}

		return nil
	})
}

func (s *Scanner) Watch(paths []string, input chan<- string, ctx context.Context) {

	s.fileWatcher.Watch(paths, input, ctx)

}
