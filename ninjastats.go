package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/assistcontrol/ninjastats/format"
	"github.com/assistcontrol/ninjastats/parser"
	"github.com/assistcontrol/ninjastats/stats"
)

func main() {
	count := flag.Bool("count", false, "Display counts")
	ci := flag.Bool("ci", false, "Display confidence intervals")
	path := flag.String("path", "/var/log/MedApps", "Path to statistics logs")
	flag.Parse()

	db := makeDB(*path)
	if *count {
		display(db, format.Headers.Count, format.Count)
	} else if *ci {
		display(db, format.Headers.CI, format.CI)
	} else {
		display(db, format.Headers.Times, format.Time)
	}
}

// makeDB does the heavy lifting. It returns a stats.DB
// fully populated with the entire stats corpus.
func makeDB(path string) *stats.DB {
	db := &stats.DB{}
	wg := sync.WaitGroup{}

	// Select all relevant files
	glob := parser.ListFiles(path)
	wg.Add(len(glob))

	// Parse each file (asynchronously)
	requests := make(chan *parser.Request, 1000)
	defer close(requests)

	for _, file := range glob {
		go func(f string) {
			parser.ParseFile(f, requests)
			wg.Done()
		}(file)
	}

	// Register each request
	go func() {
		for req := range requests {
			db.Add(req.Page, req.Verb, req.Time)
		}
	}()

	wg.Wait()
	return db
}

// display produces and displays final output as a table
func display(db *stats.DB, header string, formatter format.Formatter) {
	const tableFormat = "%15s\t%15s\t%15s\t%15s\n"

	fmt.Println(header)
	fmt.Printf(tableFormat, "URL", "GET", "POST", "HEAD")
	fmt.Println(strings.Repeat("─", 67))

	for _, page := range db.Pages() {
		fmt.Printf(tableFormat, page,
			formatter(db.NewReport(page, "GET")),
			formatter(db.NewReport(page, "POST")),
			formatter(db.NewReport(page, "HEAD")),
		)
	}
}
