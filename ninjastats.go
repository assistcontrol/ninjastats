package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// page, verb, time
var requestRE = regexp.MustCompile(`Rendered (\S+)\|(GET|POST|HEAD) in ([\d\.]+) secs`)

// request holds the data from a single request
type request = struct {
	Page string
	Verb string
	Time float64
}

var headers = struct{ count, times string }{
	count: "count (outliers):",
	times: "msec/page (pages/sec):",
}

func main() {
	count := flag.Bool("c", false, "Display counts")
	path := flag.String("path", "/var/log/MedApps", "Path to statistics logs")
	flag.Parse()

	db := makeStatsDB(*path)
	if *count {
		display(db, headers.count, countFormatter)
	} else {
		display(db, headers.times, timeFormatter)
	}
}

// makeStatsDB does the heavy lifting. It returns a StatsDB
// fully populated with the entire stats corpus.
func makeStatsDB(path string) *StatsDB {
	db := &StatsDB{}
	wg := sync.WaitGroup{}

	// Select all relevant files
	glob := ListFiles(path)
	wg.Add(len(glob))

	// Signal when all scanning is done
	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	// Parse each file (asynchronously)
	requests := make(chan *request, 1000)
	for _, file := range glob {
		go func(f string) {
			ParseFile(f, requests)
			wg.Done()
		}(file)
	}

	for {
		select {
		case req := <-requests:
			// Add each request to the DB
			db.Add(req.Page, strToVerb(req.Verb), req.Time)
		case <-done:
			// Return when done
			return db
		}
	}
}

// Type Formatter is a function that formats a single verb's Report
// as a string (< 16 chars).
type Formatter func(r *Report) string

func timeFormatter(r *Report) string {
	// If Mean == 0, rate is nonsense
	if r.Mean == 0 {
		return ""
	}

	return fmt.Sprintf("%.2f (%d)", r.Mean, r.Rate)
}

func countFormatter(r *Report) string {
	// Reduce table clutter
	if r.Count == 0 {
		return ""
	}

	return fmt.Sprintf("%d (%5d)", r.Count, r.Outliers)
}

// display produces and displays final output as a table
func display(db *StatsDB, header string, formatter Formatter) {
	tableFormat := "%15s\t%15s\t%15s\t%15s\n"

	fmt.Println(header)
	fmt.Printf(tableFormat, "URL", "GET", "POST", "HEAD")
	fmt.Println(strings.Repeat("???", 67))

	for _, page := range db.Pages() {
		fmt.Printf(tableFormat, page,
			formatter(db.NewReport(page, GET)),
			formatter(db.NewReport(page, POST)),
			formatter(db.NewReport(page, HEAD)),
		)
	}
}
