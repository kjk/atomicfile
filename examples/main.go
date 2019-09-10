package main

import (
	"fmt"
	"github.com/kjk/atomicfile"
	"log"
	"os"
)

func writeToFileAtomically(filePath string, data []byte) error {
	w, err := atomicfile.New(filePath)
	if err != nil {
		return err
	}
	// calling Close() twice is a no-op
	defer func() {
		_ = w.Close()
	}()
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return nil
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