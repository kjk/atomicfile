package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/kjk/atomicfile"
)

func writeToFileAtomically(filePath string, data []byte) error {
	f, err := atomicfile.New(filePath)
	if err != nil {
		return err
	}

	// calling Close() twice is a no-op
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return f.Close()
}

func functionThatMightPanic() {
	if rand.Intn(100) == 25 {
		panic("we crashed")
	}
}

func writeToFileAtomicallyRemoveOnPanic(filePath string, data []byte) error {
	f, err := atomicfile.New(filePath)
	if err != nil {
		return err
	}

	// RemoveIfNotClosed() after Close() is a no-op
	defer f.RemoveIfNotClosed()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	functionThatMightPanic()

	return f.Close()
}

func main() {
	fileName := "foo.txt"
	data := []byte("hello\n")
	err := writeToFileAtomically(fileName, data)
	if err != nil {
		fmt.Printf("writeToFileAtomically failed with '%s'\n", err)
		return
	}
	st, err := os.Stat(fileName)
	if err != nil {
		log.Fatalf("os.Stat('%s') failed with '%s'\n", fileName, err)
	}
	fmt.Printf("Wrote to file '%s' atomically. Size of file: %d bytes\n", fileName, st.Size())
	_ = os.Remove(fileName)
}
