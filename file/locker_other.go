//go:build !unix && !windows

package file

import (
	"errors"
	"io/fs"
)

type lockType int8

const (
	readLock = iota + 1
	writeLock
)

func lock(f File, lt lockType, immediately bool) error {
	return &fs.PathError{
		Op:   lt.String(),
		Path: f.Name(),
		Err:  errors.ErrUnsupported,
	}
}

func unlock(f File) error {
	return &fs.PathError{
		Op:   "Unlock",
		Path: f.Name(),
		Err:  errors.ErrUnsupported,
	}
}
