//go:build windows

package file

import (
	"io/fs"

	"golang.org/x/sys/windows"
)

type lockType uint32

const (
	readLock  lockType = 0
	writeLock lockType = windows.LOCKFILE_EXCLUSIVE_LOCK
)

const (
	reserved = 0
	allBytes = ^uint32(0)
)

var (
	modkernel32    = windows.NewLazySystemDLL("kernel32.dll")
	procLockFileEx = modkernel32.NewProc("LockFileEx")
)

func lock(f Locker, lt lockType, immediately bool) error {
	// Per https://golang.org/issue/19098, “Programs currently expect the Fd
	// method to return a handle that uses ordinary synchronous I/O.”
	// However, LockFileEx still requires an OVERLAPPED structure,
	// which contains the file offset of the beginning of the lock range.
	// We want to lock the entire file, so we leave the offset as zero.
	ol := new(windows.Overlapped)
	flags := uint32(lt)
	if immediately {
		flags |= windows.LOCKFILE_FAIL_IMMEDIATELY
	}
	err := windows.LockFileEx(windows.Handle(f.Fd()), flags, reserved, allBytes, allBytes, ol)
	if err != nil {
		if errno, ok := err.(windows.Errno); ok {
			if errno == windows.ERROR_LOCK_VIOLATION {
				err = errAlreadyLocked
			}
		}

		return &fs.PathError{
			Op:   lt.String(),
			Path: f.Name(),
			Err:  err,
		}
	}
	return nil
}

func unlock(f Locker) error {
	ol := new(windows.Overlapped)
	err := windows.UnlockFileEx(windows.Handle(f.Fd()), reserved, allBytes, allBytes, ol)
	if err != nil {
		return &fs.PathError{
			Op:   "Unlock",
			Path: f.Name(),
			Err:  err,
		}
	}
	return nil
}
