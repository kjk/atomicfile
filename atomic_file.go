package atomicfile

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Writer allows writing to a file atomically
// i.e. if the while file is not written successfully, we make sure
// to clean things up
type Writer struct {
	dstPath string
	tmpFile *os.File
	err     error

	tmpPath string // for debugging
}

// NewWriter creates new Writer
func NewWriter(path string) (*Writer, error) {
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

	return &Writer{
		dstPath: path,
		tmpFile: tmpFile,
		tmpPath: tmpFile.Name(),
	}, nil
}

// Write writes data to a file
func (w *Writer) Write(d []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	n, err := w.tmpFile.Write(d)
	if err != nil {
		w.err = err
		// cleanup i.e. delete temporary file
		w.Close()
		return 0, err
	}
	return n, nil
}

// Cancel cancels writing and removes the temp file.
// Destination file will not be created
// Use it to cleanup things when error happens outside of Write()
// Cancel after Close is harmless to make it easier to use via defer
func (w *Writer) Cancel() {
	w.err = errors.New("cancelled")
	w.Close()
}

// Close closes the file. Can be called multiple times to make it
// easier to use via defer
func (w *Writer) Close() error {
	if w.tmpFile == nil {
		// was already called, return same error as first Close()
		return w.err
	}
	tmpFile := w.tmpFile
	w.tmpFile = nil

	// cleanup things (delete temporary files) if:
	// - there was an error in Write()
	// - Close() failed
	// - rename to destination failed
	err := tmpFile.Close()

	defer os.Remove(w.tmpPath) // ignoring error on this one

	// if there was an error during write, return that error
	if w.err != nil {
		return w.err
	}

	if err == nil {
		// this might over-write dstPath
		err = os.Rename(w.tmpPath, w.dstPath)
	}
	w.err = err
	return w.err
}
