package main

import (
	"flag"
	"fmt"
	"regexp"
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
	var (
		db       = &StatsDB{}
		wg       = sync.WaitGroup{}
		requests = make(chan *request, 1000)
		done     = make(chan bool)
	)

	// Select all relevant files
	glob := listFiles(path)
	wg.Add(len(glob))

	// Signal when all scanning is done
	go func() {
		wg.Wait()
		done <- true
	}()

	// Parse each file (asynchronously)
	for _, file := range glob {
		go func(f string) {
			defer func() { wg.Done() }()
			ParseFile(f, requests)
		}(file)
	}

	for {
		select {
		// Add each request to the DB
		case req := <-requests:
			verb := strToVerb(req.Verb)
			if verb == __FAILED__ {
				continue
			}
			db.Add(req.Page, verb, req.Time)
			// Return when done
		case <-done:
			return db
		}
	}
}

func timeFormatter(r *VerbReport) string {
	// If Mean == 0, rate is nonsense
	if r.Mean == 0 {
		return ""
	}

	return fmt.Sprintf("%.2f (%d)", r.Mean, r.Rate)
}

func countFormatter(r *VerbReport) string {
	// Reduce table clutter
	if r.Count == 0 {
		return ""
	}

	return fmt.Sprintf("%d (%4d)", r.Count, r.Outliers)
}

// display produces and displays final output as a table
func display(db *StatsDB, header string, formatter func(*VerbReport) string) {
	tableFormat := "%15s\t%15s\t%15s\t%15s\n"

	fmt.Println(header)
	fmt.Printf(tableFormat, "URL", "GET", "POST", "HEAD")
	fmt.Println("───────────────────────────────────────────────────────────────────")

	report := db.Report()
	for _, p := range db.Pages() {
		page := report[p]

		fmt.Printf(tableFormat, p, formatter(page.Get), formatter(page.Post), formatter(page.Head))
	}
}
