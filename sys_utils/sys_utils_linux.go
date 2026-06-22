package sysutils

import (
	"os"
	"syscall"
)

func LockFile(file *os.File) error {
	return syscall.Flock(int(file.Fd()), syscall.LOCK_SH)

}

func UnlockFile(file *os.File) error {
	return syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
}
