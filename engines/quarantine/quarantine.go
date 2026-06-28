package quarantine

import (
	"SimpleAV/apperrors"
	"SimpleAV/config"
	sysutils "SimpleAV/sys_utils"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type Quarantiner struct{}

func NewQuarantiner() *Quarantiner {
	return &Quarantiner{}
}

type QuarantineMeta struct {
	ID              string
	OriginalPath    string
	QuarantinedAt   time.Time
	DetectionMethod string
	SHA256          string
	FileSize        int64
	Status          string
}

func Init() error {
	err := sysutils.EnsureDir(config.QuarantineDir)
	if err != nil {
		fmt.Println("quarantine directory creaton is affected %w", err)
		return err
	}
	return nil
}

func (q *Quarantiner) Quarantine(filePath string) error {

	// check if file exists and can be read

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("%w:%w", apperrors.ErrFileOpen, err)
	}
	defer f.Close()

	// locks to avoid TOCTOU race
	sysutils.LockFile(f)
	defer sysutils.UnlockFile(f)

	info, err := f.Stat()

	// hash the original file
	hash, err := sysutils.ConvertToSHA256FromFile(f)
	if err != nil {
		return fmt.Errorf("hashing err %w:%w", apperrors.ErrHashing, err)
	}

	id := uuid.New().String()
	quarPath := filepath.Join(config.QuarantineDir, id+".quar")
	metaPath := filepath.Join(config.QuarantineDir, id+".meta.json")

	dst, err := os.OpenFile(
		quarPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0600,
	)
	if err != nil {
		return fmt.Errorf("%w:%w", apperrors.ErrFileOpen, err)
	}
	defer dst.Close()

	// reset ptr position to start(0)
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// malform and write mangled file to quarantine dir
	if _, err := io.Copy(dst, NewXorReader(f, 0xAA)); err != nil {
		return fmt.Errorf("%w:%w", apperrors.ErrDataCopy, err)
	}

	// write metadata
	meta := QuarantineMeta{
		ID:            id,
		OriginalPath:  filePath,
		QuarantinedAt: time.Now().UTC(),
		SHA256:        hash,
		FileSize:      info.Size(),
		Status:        "quarantined",
	}

	metaBytes, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(metaPath, metaBytes, 0600); err != nil {
		return err
	}

	// delete the original file
	return removeOriginal(filePath)
}

type xorReader struct {
	r   io.Reader
	key byte
}

func NewXorReader(r io.Reader, key byte) io.Reader {
	return &xorReader{
		r:   r,
		key: key,
	}
}

func (x *xorReader) Read(p []byte) (int, error) {
	n, err := x.r.Read(p)

	for i := 0; i < n; i++ {
		p[i] ^= x.key
	}

	return n, err
}
