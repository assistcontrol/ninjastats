package parser

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/assistcontrol/ninjastats/parser/bz2"
	"github.com/assistcontrol/ninjastats/parser/flat"
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
	var file *os.File
	var scanner *bufio.Scanner

	if strings.HasSuffix(path, ".bz2") {
		file = bz2.Open(path)
		scanner = bz2.Scanner(file)
	} else {
		file = flat.Open(path)
		scanner = flat.Scanner(file)
	}
	defer file.Close()

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
