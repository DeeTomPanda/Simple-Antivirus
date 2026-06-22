package sysutils

import (
	"os"

	"golang.org/x/sys/windows"
)

func LockFile(file *os.File) error {
	var ol windows.Overlapped

	return windows.LockFileEx(
		windows.Handle(file.Fd()),
		windows.LOCKFILE_EXCLUSIVE_LOCK,
		0,
		1,
		0,
		&ol,
	)

}

func UnlockFile(file *os.File) error {
	var ol windows.Overlapped

	return windows.UnlockFileEx(
		windows.Handle(file.Fd()),
		0,
		1,
		0,
		&ol)
}
