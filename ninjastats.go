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

var headers = struct{ ci, count, times string }{
	ci:    "msec/page ± 95% CI:",
	count: "count (outliers):",
	times: "msec/page (pages/sec):",
}

func main() {
	count := flag.Bool("c", false, "Display counts")
	ci := flag.Bool("ci", false, "Display confidence intervals")
	path := flag.String("path", "/var/log/MedApps", "Path to statistics logs")
	flag.Parse()

	db := makeStatsDB(*path)
	if *count {
		display(db, headers.count, countFormatter)
	} else if *ci {
		display(db, headers.ci, ciFormatter)
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

	// Parse each file (asynchronously)
	requests := make(chan *request, 1000)
	for _, file := range glob {
		go func(f string) {
			ParseFile(f, requests)
			wg.Done()
		}(file)
	}

	// Register each request
	go func() {
		for req := range requests {
			db.Add(req.Page, strToVerb(req.Verb), req.Time)
		}
	}()

	wg.Wait()
	return db
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

func ciFormatter(r *Report) string {
	if r.Mean == 0 {
		return ""
	}

	return fmt.Sprintf("%.2f ± %.2f", r.Mean, r.CIWidth)
}

// display produces and displays final output as a table
func display(db *StatsDB, header string, formatter Formatter) {
	tableFormat := "%15s\t%15s\t%15s\t%15s\n"

	fmt.Println(header)
	fmt.Printf(tableFormat, "URL", "GET", "POST", "HEAD")
	fmt.Println(strings.Repeat("─", 67))

	for _, page := range db.Pages() {
		fmt.Printf(tableFormat, page,
			formatter(db.NewReport(page, GET)),
			formatter(db.NewReport(page, POST)),
			formatter(db.NewReport(page, HEAD)),
		)
	}
}
