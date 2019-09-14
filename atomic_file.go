package atomicfile

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	// ErrCancelled is returned by calls subsequent to Cancel()
	ErrCancelled = errors.New("cancelled")

	// ensure we implement desired interface
	_ io.WriteCloser = &File{}
)

// File allows writing to a file atomically
// i.e. if the while file is not written successfully, we make sure
// to clean things up
type File struct {
	dstPath string
	tmpFile *os.File
	err     error

	tmpPath string // for debugging
}

// New creates new File
func New(path string) (*File, error) {
	base, fName := filepath.Split(path)
	if base == "" {
		base = "."
	}
	if fName == "" {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrInvalid}
	}

	tmpFile, err := ioutil.TempFile(base, fName)
	if err != nil {
		return nil, err
	}

	return &File{
		dstPath: path,
		tmpFile: tmpFile,
		tmpPath: tmpFile.Name(),
	}, nil
}

// Write writes data to a file
func (w *File) Write(d []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	n, err := w.tmpFile.Write(d)
	if err != nil {
		w.err = err
		// cleanup i.e. delete temporary file
		_ = w.Close()
		return 0, err
	}
	return n, nil
}

// Cancel cancels writing and removes the temp file.
// Destination file will not be created
// Use it to cleanup things when error happens outside of Write()
// Cancel after Close is harmless to make it easier to use via defer
func (w *File) Cancel() {
	w.err = ErrCancelled
	_ = w.Close()
}

// Close closes the file. Can be called multiple times to make it
// easier to use via defer
func (w *File) Close() error {
	if w.tmpFile == nil {
		// was already called, return same error as first Close()
		return w.err
	}
	tmpFile := w.tmpFile
	w.tmpFile = nil

	// cleanup things (delete temporary files) if:
	// - there was an error in Write()
	// - thre was an error in Sync()
	// - Close() failed
	// - rename to destination failed

	// https://www.joeshaw.org/dont-defer-close-on-writable-files/
	errSync := tmpFile.Sync()
	errClose := tmpFile.Close()

	// always delete the temporary file
	defer func() {
		// ignoring error on this one
		_ = os.Remove(w.tmpPath)
	}()

	// if there was an error during write, return that error
	if w.err != nil {
		return w.err
	}

	err := errSync
	if err == nil {
		err = errClose
	}

	if err == nil {
		// this will over-write dstPath (if it exists)
		err = os.Rename(w.tmpPath, w.dstPath)
	}
	if w.err == nil {
		w.err = err
	}
	return w.err
}
