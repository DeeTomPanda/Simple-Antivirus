package quarantine

import (
	"SimpleAV/config"

	"bytes"
	"os"
	"path/filepath"
	"testing"
)

type mockQuarantiner struct{}

func (m *mockQuarantiner) Quarantine(filePath string) error {
	return nil
}

var meta QuarantineMeta

func createFileWithData(t *testing.T, data []byte) string {
	t.Helper()

	dir := t.TempDir()

	p := filepath.Join(dir, "bad.exe")

	err := os.WriteFile(p, data, 0644)
	if err != nil {
		t.Fatal(err)
	}

	return p
}

func TestQuarantineHappyPath(t *testing.T) {

	config.QuarantineDir = t.TempDir()

	file := createFileWithData(
		t,
		[]byte("malicious payload"),
	)

	q := NewQuarantiner()

	err := q.Quarantine(file)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(file)

	if !os.IsNotExist(err) {
		t.Fatal("expected original file removed")
	}

	entries, err := os.ReadDir(config.QuarantineDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 2 {
		t.Fatalf(
			"expected 2 files, got %d",
			len(entries),
		)
	}
}

func TestQuarantineZeroByteFile(t *testing.T) {

	config.QuarantineDir = t.TempDir()

	file := createFileWithData(t, nil)

	q := NewQuarantiner()

	err := q.Quarantine(file)

	if err != nil {
		t.Fatal(err)
	}

	entries, _ := os.ReadDir(config.QuarantineDir)

	if len(entries) != 2 {
		t.Fatal("expected quarantine artifacts")
	}
}

func TestQuarantineMissingFile(t *testing.T) {

	config.QuarantineDir = t.TempDir()

	q := NewQuarantiner()

	err := q.Quarantine("missing.exe")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestXORRoundTrip(t *testing.T) {

	original := []byte("evil")

	mangled := make([]byte, len(original))
	copy(mangled, original)

	for i := range mangled {
		mangled[i] ^= 0xAA
	}

	for i := range mangled {
		mangled[i] ^= 0xAA
	}

	if !bytes.Equal(original, mangled) {
		t.Fatal("xor not reversible")
	}
}
