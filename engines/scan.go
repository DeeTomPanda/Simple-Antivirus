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

type Scanner struct {
	checker HashChecker
}

func NewScanner(checker HashChecker) *Scanner {
	return &Scanner{checker: checker}
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

		malicious, err := s.checker.CheckMaliciousHash(path)
		if err != nil {
			return err
		}

		if malicious {
			applogger.Warn("malware detected: " + path)
		}

		return nil
	})
}
