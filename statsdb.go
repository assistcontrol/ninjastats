package main

import (
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
	StatsDB    map[string]*Page
	Report     map[string]*PageReport
	VerbReport struct {
		Mean     float64
		Rate     int
		Count    int
		Outliers int
	}
	PageReport struct{ Get, Post, Head *VerbReport }
)

// Add is the main entry function that registers a time for a
// given page and HTTP verb.
func (sdb StatsDB) Add(page string, verb Verb, tm float64) {
	if _, exists := sdb[page]; !exists {
		sdb[page] = NewPage()
	}

	sdb.times(page, verb).Add(tm)
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

func (sdb StatsDB) Report() Report {
	r := make(Report)
	for _, page := range sdb.Pages() {
		r[page] = sdb.pageReport(page)
	}

	return r
}

// times returns the underlying Times structure for a given page
// and verb
func (sdb StatsDB) times(page string, verb Verb) *Times {
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

// verbReport returns a VerbReport structure detailing a single verb from
// a single page
func (sdb StatsDB) verbReport(page string, verb Verb) *VerbReport {
	base := sdb.times(page, verb)
	reduced := base.Reduce()
	mean := reduced.Mean()

	return &VerbReport{
		Mean:     math.Round(mean*1000*100) / 100, // s -> ms, round to 2 places
		Rate:     int(math.Round(1.0 / mean)),
		Count:    base.Count(),
		Outliers: base.Count() - reduced.Count(),
	}
}

// pageReport collects VerbReports for a single page
func (sdb StatsDB) pageReport(page string) *PageReport {
	return &PageReport{
		Get:  sdb.verbReport(page, GET),
		Post: sdb.verbReport(page, POST),
		Head: sdb.verbReport(page, HEAD),
	}
}
