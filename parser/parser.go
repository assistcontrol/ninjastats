package parser

import (
	"bufio"
	"compress/bzip2"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// page, verb, time
var requestRE = regexp.MustCompile(`Rendered (\S+)\|(GET|POST|HEAD) in ([\d\.]+) secs`)

// Request holds the data from a single request
type Request struct {
	Page string
	Verb string
	Time float64
}

// listFiles returns a slice of paths to logs in the given dir
func ListFiles(path string) []string {
	glob, err := filepath.Glob(path + "/statistics.log*")
	if err != nil {
		log.Fatal("glob:", err)
	}
	if len(glob) == 0 {
		log.Fatalf("%s/statistics.log*: No files found", path)
	}

	return glob
}

// ParseFile does the heavy lifting of extracting requests
// from a given file. Each extracted request is sent up
// the supplied channel.
func ParseFile(path string, reqChan chan<- *Request) {
	compressed := strings.HasSuffix(path, ".bz2")

	// It's weird to open the file here, but to defer Close() it
	// has to appear in the same function as the scanner
	file := openFile(path)
	defer file.Close()
	scanner := newScanner(newReader(file, compressed))

	for scanner.Scan() {
		req, ok := parseLine(scanner.Text())
		if !ok {
			continue
		}

		reqChan <- req
	}
}

// openFile returns a handle for an opened log file
func openFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}

	return file
}

// newReader returns a reader capable of reading the file
// argument. If compressed is true, the file is assumed to be
// bzip2-compressed and is passed through an appropriate
// decompressor.
func newReader(file *os.File, compressed bool) io.Reader {
	if compressed {
		return bzip2.NewReader(file)
	}

	return file
}

// newScanner returns a line-oriented bufio.Scanner for a
// given reader
func newScanner(reader io.Reader) *bufio.Scanner {
	return bufio.NewScanner(reader)
}

// parseLine uses the requestRE regexp to extract request
// data. Returns a *request struct and a bool indicating
// whether the extraction was successful.
func parseLine(s string) (*Request, bool) {
	// Match [entire string, page, verb, time]
	match := requestRE.FindStringSubmatch(s)
	if match == nil {
		// Did not match
		return nil, false
	}

	tm, err := strconv.ParseFloat(match[3], 64)
	if err != nil {
		// Badly-formatted time
		return nil, false
	}

	// Matched successfully
	return &Request{
		Page: match[1],
		Verb: match[2],
		Time: tm,
	}, true
}
