package bz2

// bz2 is a simple wrapper around the bzip2 package that provides a
// convenient way to open a bzip2 file and get a scanner for it.

import (
	"bufio"
	"compress/bzip2"
	"io"
	"log"
	"os"
)

// Open opens a bzip2 file and returns a file handle.
func Open(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}

	return file
}

// Scanner returns a scanner for a bzip2 file.
func Scanner(file io.Reader) *bufio.Scanner {
	return bufio.NewScanner(bzip2.NewReader(file))
}
