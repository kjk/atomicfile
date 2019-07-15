package atomicfile

import (
	"io/ioutil"
	"os"
	"testing"
)

func assertFileExists(t *testing.T, path string) {
	st, err := os.Stat(path)
	if err != nil {
		t.Fatalf("file '%s' doesn't exist, os.Stat() failed with '%s'", path, err)
	}
	if !st.Mode().IsRegular() {
		t.Fatalf("Path '%s' exists but is not a file (mode: %d)", path, int(st.Mode()))
	}
}

func assertFileNotExists(t *testing.T, path string) {
	_, err := os.Stat(path)
	if err == nil {
		t.Fatalf("file '%s' exist, expected to not exist", path)
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("error: %s", err)
	}
}

func assertError(t *testing.T, err error) {
	if err == nil {
		t.Fatal("expected to get an error")
	}
}

func assertFileSizeEqual(t *testing.T, path string, n int64) {
	st, err := os.Stat(path)
	if err != nil {
		t.Fatalf("os.Stat('%s') failed with '%s'", path, err)
	}
	if st.Size() != n {
		t.Fatalf("path: '%s', expected size: %d, got: %d", path, n, st.Size())
	}
}

func assertIntEqual(t *testing.T, exp int, got int) {
	if exp != got {
		t.Fatalf("expected: %d, got: %d", exp, got)
	}
}
func TestWrite(t *testing.T) {
	dst := "atomic_file.go.copy"
	os.Remove(dst)
	{
		w, err := NewWriter(dst)
		assertNoError(t, err)
		assertFileExists(t, w.tmpPath)
		w.Close()
		assertFileExists(t, dst)
		assertFileSizeEqual(t, dst, 0)
		assertFileNotExists(t, w.tmpPath)
	}
	d, err := ioutil.ReadFile("atomic_file.go")
	assertNoError(t, err)

	{
		w, err := NewWriter(dst)
		assertNoError(t, err)
		assertFileExists(t, w.tmpPath)
		n, err := w.Write(d)
		assertNoError(t, err)
		assertIntEqual(t, n, len(d))
		assertFileExists(t, w.tmpPath)
		err = w.Close()
		assertNoError(t, err)
		assertFileNotExists(t, w.tmpPath)
		assertFileSizeEqual(t, dst, int64(len(d)))
		// calling Close twice is a no-op
		err = w.Close()
		assertNoError(t, err)
	}
	os.Remove(dst)
}
