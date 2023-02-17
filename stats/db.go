package stats

import (
	"log"
	"math"
	"sort"
)

// Verb
type Verb int

const (
	GET Verb = iota
	POST
	HEAD
	__FAILED__
)

// strToVerb converts a string to a Verb enum value
func strToVerb(verb string) Verb {
	switch verb {
	case "GET":
		return GET
	case "POST":
		return POST
	case "HEAD":
		return HEAD
	}

	log.Fatal("Unknown verb:", verb)
	return __FAILED__
}

// page
type page struct{ get, post, head *times }

func newPage() *page {
	return &page{
		get:  newTimes(),
		post: newTimes(),
		head: newTimes(),
	}
}

// DB
type (
	DB     map[string]*page
	Report struct {
		Mean     float64
		Rate     int
		Count    int
		Outliers int
		CIWidth  float64
	}
)

// Add is the main entry function that registers a time for a
// given page and HTTP verb.
func (sdb DB) Add(page, verb string, tm float64) {
	if _, exists := sdb[page]; !exists {
		sdb[page] = newPage()
	}

	sdb.getTimes(page, verb).add(tm)
}

// Pages returns a sorted list of known pages.
func (sdb DB) Pages() []string {
	pages := make([]string, 0, len(sdb))
	for pg := range sdb {
		pages = append(pages, pg)
	}
	sort.Strings(pages)

	return pages
}

// NewReport returns a Report structure detailing a single verb from
// a single page
func (sdb DB) NewReport(page, verb string) *Report {
	base := sdb.getTimes(page, verb)
	reduced := base.reduce()
	mean := reduced.mean()
	width := reduced.ciWidth()

	return &Report{
		Mean:     secToMsec(mean),
		Rate:     int(math.Round(1.0 / mean)),
		Count:    base.count(),
		Outliers: base.count() - reduced.count(),
		CIWidth:  secToMsec(width),
	}
}

// times returns the underlying Times structure for a given page
// and verb
func (sdb DB) getTimes(page, verb string) *times {
	switch strToVerb(verb) {
	case GET:
		return sdb[page].get
	case POST:
		return sdb[page].post
	case HEAD:
		return sdb[page].head
	}

	return nil
}

// secToMsec converts a time in seconds to milliseconds, rounded to
// two decimal places
func secToMsec(sec float64) float64 {
	return math.Round(sec*1000*100) / 100
}
