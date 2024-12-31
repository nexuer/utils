package file

import (
	"errors"
	"io/fs"
)

var (
	errAlreadyLocked = errors.New("file already locked")
)

type Locker interface {
	// Name returns the name of the file.
	Name() string

	// Fd returns a valid file descriptor.
	// (If the File is an *os.File, it must not be closed.)
	Fd() uintptr

	// Stat returns the FileInfo structure describing file.
	Stat() (fs.FileInfo, error)
}

func IsAlreadyLocked(err error) bool {
	var e *fs.PathError
	ok := errors.As(err, &e)
	return ok && errors.Is(e.Err, errAlreadyLocked)
}

// String returns the name of the function corresponding to lt
// (Lock, RLock, or Unlock).
func (lt lockType) String() string {
	switch lt {
	case readLock:
		return "RLock"
	case writeLock:
		return "Lock"
	default:
		return "Unlock"
	}
}

// Lock places an advisory write lock on the file, When the immediately is false, it will block
// until the lock can be locked, otherwise an error will be returned immediately.
//
// If Lock returns nil, no other process will be able to place a read or write
// lock on the file until this process exits, closes f, or calls Unlock on it.
//
// If f's descriptor is already read- or write-locked, the behavior of Lock is
// unspecified.
//
// Closing the file may or may not release the lock promptly. Callers should
// ensure that Unlock is always called when Lock succeeds.
func Lock(f Locker, immediately bool) error {
	return lock(f, writeLock, immediately)
}

// RLock places an advisory read lock on the file, When the immediately is false, it will block
// until the lock can be locked, otherwise an error will be returned immediately.
//
// If RLock returns nil, no other process will be able to place a write lock on
// the file until this process exits, closes f, or calls Unlock on it.
//
// If f is already read- or write-locked, the behavior of RLock is unspecified.
//
// Closing the file may or may not release the lock promptly. Callers should
// ensure that Unlock is always called if RLock succeeds.
func RLock(f Locker, immediately bool) error {
	return lock(f, readLock, immediately)
}

// Unlock removes an advisory lock placed on f by this process.
//
// The caller must not attempt to unlock a file that is not locked.
func Unlock(f Locker) error {
	return unlock(f)
}
