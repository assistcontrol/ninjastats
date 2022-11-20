package main

import "math"

const StdDevs = 3 // How many deviations from the mean to accept

// Times
type Times struct {
	times []float64
}

// NewTimes creates the Times struct. Required for initialization of
// the underlying .times slice.
func NewTimes() *Times {
	return &Times{
		times: make([]float64, 0),
	}
}

// Add appends a time tm to the collection of times
func (ts *Times) Add(tm float64) {
	ts.times = append(ts.times, tm)
}

// Count returns the total number of times (including outliers)
func (ts Times) Count() int {
	return len(ts.times)
}

// Mean returns the mean time, or 0 for an empty list
func (ts Times) Mean() float64 {
	N := ts.Count()
	if N == 0 {
		return 0
	}

	return ts.sum() / float64(ts.Count())
}

// Reduce produces a new Times struct from the parent data with
// outliers removed
func (ts Times) Reduce() *Times {
	mean := ts.Mean()
	limit := StdDevs * ts.stddev()

	reduced := &Times{}
	for _, t := range ts.times {
		if math.Abs(t-mean) <= limit {
			reduced.Add(t)
		}
	}

	return reduced
}

// stddev returns the population standard deviation of the times
func (ts Times) stddev() float64 {
	N := float64(ts.Count())
	if N == 0 {
		return 0
	}

	return math.Sqrt(ts.variance() / N)
}

// sum returns the sum of all times
func (ts Times) sum() (sum float64) {
	for _, t := range ts.times {
		sum += t
	}

	return
}

// variance returns the population variance of the times
func (ts Times) variance() (v float64) {
	mean := ts.Mean()
	for _, t := range ts.times {
		v += (t - mean) * (t - mean)
	}

	return
}
