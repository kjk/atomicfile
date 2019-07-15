# atomicfile

In Go, when writing to a file, we should:

- handle error returned by `Close()` method
- remove partially written file if `Write()` or `Close()` return an error

This logic is needed often and non-trivial to write.

Package `atomicfile` encapsulates this logic and also ensure atomic writing
to a file.

We write to a temporary file and only if it was successfully written, we rename
it atomically to a desired final name.

The usage is:

```go
w, err := atomicfile.NewWriter("foo.txt")
if err != nil {
    // handle error
    log.Fatalf("atomicfile.NewWriter() failed with %s\n", err)
}
// calling Close() twice is a no-op
defer w.Close()
_, err := w.Write([]byte("hello"))
if err != nil {
    // handle error
    log.Fatalf("w.Writer() failed with %s\n", err)
}
err = w.Close()
if err != nil {
    // handle error
    log.Fatalf("w.Close() failed with %s\n", err)
}
// we've atomically created foo.txt
```
