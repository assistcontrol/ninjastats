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
	path := flag.String("path", "/var/log/ninja", "Path to statistics logs")
	flag.Parse()

	db := makeDB(*path)

	tables := []string{
		display(db, format.HeaderTimes, format.Time),
		display(db, format.HeaderCount, format.Count),
		display(db, format.HeaderCI, format.CI),
	}
	fmt.Print(strings.Join(tables, "\n"))
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
	requests := make(chan *parser.Request, 10)
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
func display(db *stats.DB, header string, formatter format.Formatter) string {
	const tableFormat = "%15s\t%15s\t%15s\t%15s\n"

	s := fmt.Sprintln(header)
	s += fmt.Sprintf(tableFormat, "URL", "GET", "POST", "HEAD")
	s += fmt.Sprintln(strings.Repeat("─", 67))

	for _, page := range db.Pages() {
		s += fmt.Sprintf(tableFormat, page,
			formatter(db.NewReport(page, "GET")),
			formatter(db.NewReport(page, "POST")),
			formatter(db.NewReport(page, "HEAD")),
		)
	}

	return s
}
