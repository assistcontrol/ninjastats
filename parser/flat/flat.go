package flat

// flat is a simple wrapper around the standard library that provides a
// convenient way to open a file and get a scanner for it.

import (
	"bufio"
	"io"
	"log"
	"os"
)

// Open opens a file and returns a file handle.
func Open(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}

	return file
}

// Scanner returns a scanner for a file.
func Scanner(file io.Reader) *bufio.Scanner {
	return bufio.NewScanner(file)
}
