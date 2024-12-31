//go:build unix

package file

import (
	"io/fs"

	"golang.org/x/sys/unix"
)

type lockType int16

const (
	readLock  lockType = unix.LOCK_SH
	writeLock lockType = unix.LOCK_EX
)

func lock(f Locker, lt lockType, immediately bool) error {
	flags := lt
	if immediately {
		flags |= unix.LOCK_NB
	}

	if err := unix.Flock(int(f.Fd()), int(flags)); err != nil {
		if errno, ok := err.(unix.Errno); ok {
			if errno == unix.EWOULDBLOCK {
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
	return lock(f, unix.LOCK_UN, false)
}
