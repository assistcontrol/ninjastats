package parser

import (
	"bufio"
	"compress/bzip2"
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
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic("Unable to close " + path + ": " + err.Error())
		}
	}()

	// Create a scanner for the file
	var scanner *bufio.Scanner
	switch {
	case strings.HasSuffix(path, ".bz2"):
		scanner = bufio.NewScanner(bzip2.NewReader(file))
	default:
		scanner = bufio.NewScanner(file)
	}

	for scanner.Scan() {
		req, ok := parseLine(scanner.Text())
		if !ok {
			continue
		}

		reqChan <- req
	}
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
