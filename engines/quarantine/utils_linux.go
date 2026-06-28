package quarantine

import (
	"SimpleAV/apperrors"
	"fmt"
	"os"
)

func removeOriginal(path string) error {
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("%w:%w", apperrors.ErrFileDeletion, err)
	}
	return nil
}
