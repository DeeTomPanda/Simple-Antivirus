//go:build windows

package quarantine

import (
	"SimpleAV/apperrors"
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/windows"
)

func removeOriginal(path string) error {
	if err := os.Remove(path); err == nil {
		return nil
	}
	// retry with backoff
	if err := removeWithRetry(path, 5); err == nil {
		return nil
	}
	// last resort. delete on reboot
	return markForDeletion(path)
}

func removeWithRetry(path string, attempts int) error {
	var err error
	for i := range attempts {
		err = os.Remove(path)
		if err == nil {
			return nil
		}
		// exponential backoff strategy
		time.Sleep(time.Duration(i*500) * time.Millisecond)
	}
	return fmt.Errorf("failed to remove after %d attempts: %w:%w", attempts, errors.ErrUnsupported, err)
}

func markForDeletion(path string) error {
	// MoveFileEx with MOVEFILE_DELAY_UNTIL_REBOOT
	// Windows will delete it on next boot
	// https://learn.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-movefileexa#:~:text=MOVEFILE%5FDELAY%5FUNTIL%5FREBOOT,-4
	from, _ := windows.UTF16PtrFromString(path)
	err := windows.MoveFileEx(from, nil, windows.MOVEFILE_DELAY_UNTIL_REBOOT)
	if err != nil {
		return fmt.Errorf("%w:%w", apperrors.MarkErrFileDeletion, err)
	}
	return nil
}
