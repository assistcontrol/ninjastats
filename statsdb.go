package main

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

// Page
type Page struct{ get, post, head *Times }

func NewPage() *Page {
	return &Page{
		get:  NewTimes(),
		post: NewTimes(),
		head: NewTimes(),
	}
}

// StatsDB
type (
	StatsDB map[string]*Page
	Report  struct {
		Mean     float64
		Rate     int
		Count    int
		Outliers int
		CIWidth  float64
	}
)

// Add is the main entry function that registers a time for a
// given page and HTTP verb.
func (sdb StatsDB) Add(page string, verb Verb, tm float64) {
	if _, exists := sdb[page]; !exists {
		sdb[page] = NewPage()
	}

	sdb.getTimes(page, verb).Add(tm)
}

// Pages returns a sorted list of known pages.
func (sdb StatsDB) Pages() []string {
	pages := make([]string, 0, len(sdb))
	for pg := range sdb {
		pages = append(pages, pg)
	}
	sort.Strings(pages)

	return pages
}

// times returns the underlying Times structure for a given page
// and verb
func (sdb StatsDB) getTimes(page string, verb Verb) *Times {
	switch verb {
	case GET:
		return sdb[page].get
	case POST:
		return sdb[page].post
	case HEAD:
		return sdb[page].head
	}

	return nil
}

// NewReport returns a Report structure detailing a single verb from
// a single page
func (sdb StatsDB) NewReport(page string, verb Verb) *Report {
	base := sdb.getTimes(page, verb)
	reduced := base.Reduce()
	mean := reduced.Mean()
	width := reduced.CIWidth()

	return &Report{
		Mean:     secToMsec(mean),
		Rate:     int(math.Round(1.0 / mean)),
		Count:    base.Count(),
		Outliers: base.Count() - reduced.Count(),
		CIWidth:  secToMsec(width),
	}
}

// secToMsec converts a time in seconds to milliseconds, rounded to
// two decimal places
func secToMsec(sec float64) float64 {
	return math.Round(sec*1000*100) / 100
}
