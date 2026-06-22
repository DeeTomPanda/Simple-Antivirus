package engines_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"SimpleAV/apperrors"
	"SimpleAV/engines"
)

// mockChecker test stub
type mockChecker struct {
	malicious map[string]bool
	err       error
}

func (m *mockChecker) CheckMaliciousHash(path string) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.malicious[path], nil
}

// helper for temp files
func createTempFiles(t *testing.T, names ...string) (string, []string) {
	t.Helper()
	dir := t.TempDir()
	var paths []string
	for _, name := range names {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte("data"), 0644); err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		paths = append(paths, p)
	}
	return dir, paths
}

// Case 1: malware detected
func TestScanner_MalwareDetected(t *testing.T) {
	dir, paths := createTempFiles(t, "bad.exe", "good.txt")

	checker := &mockChecker{
		malicious: map[string]bool{
			paths[0]: true,
			paths[1]: false,
		},
	}

	scanner := engines.NewScanner(checker)
	err := scanner.ScanDirectory(dir, context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Case 2: all files clean
func TestScanner_CleanFiles(t *testing.T) {
	dir, paths := createTempFiles(t, "clean1.txt", "clean2.txt")

	checker := &mockChecker{
		malicious: map[string]bool{
			paths[0]: false,
			paths[1]: false,
		},
	}

	scanner := engines.NewScanner(checker)
	err := scanner.ScanDirectory(dir, context.Background())
	if err != nil {
		t.Fatalf("expected no error on clean files, got: %v", err)
	}
}

// Case 3: context cancelled before scan starts
func TestScanner_ContextCancelledBeforeScan(t *testing.T) {
	dir, _ := createTempFiles(t, "file.txt")

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately before scan

	scanner := engines.NewScanner(&mockChecker{})
	err := scanner.ScanDirectory(dir, ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got: %v", err)
	}
}

// cancelOnCallChecker, cancel after N scans
type cancelOnCallChecker struct {
	malicious map[string]bool
	cancel    context.CancelFunc
	cancelAt  int //  count of files to be scanned before cancellation
	calls     *int
}

func (c *cancelOnCallChecker) CheckMaliciousHash(path string) (bool, error) {
	*c.calls++
	if *c.calls >= c.cancelAt {
		c.cancel()
	}
	return c.malicious[path], nil
}

// Case 4: context cancelled mid scan
func TestScanner_ContextCancelledMidScan(t *testing.T) {
	dir, paths := createTempFiles(t, "file1.txt", "file2.txt", "file3.txt")

	ctx, cancel := context.WithCancel(context.Background())

	callCount := 0
	// cancel after first file check
	cancellingChecker := &cancelOnCallChecker{
		cancel:    cancel,
		cancelAt:  1,
		malicious: map[string]bool{paths[0]: false, paths[1]: false, paths[2]: false},
		calls:     &callCount,
	}

	scanner := engines.NewScanner(cancellingChecker)
	err := scanner.ScanDirectory(dir, ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled mid-scan, got: %v", err)
	}
}

// Case 5: checker returns an error (e.g. DB down)
func TestScanner_CheckerError(t *testing.T) {
	dir, _ := createTempFiles(t, "file.txt")

	dbErr := apperrors.ErrDatabaseDown
	checker := &mockChecker{err: dbErr}

	scanner := engines.NewScanner(checker)
	err := scanner.ScanDirectory(dir, context.Background())
	if !errors.Is(err, apperrors.ErrDatabaseDown) {
		t.Fatalf("expected db error to propagate, got: %v", err)
	}
}

// Case 6: empty directory
func TestScanner_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	scanner := engines.NewScanner(&mockChecker{})
	err := scanner.ScanDirectory(dir, context.Background())
	if err != nil {
		t.Fatalf("expected no error on empty dir, got: %v", err)
	}
}
