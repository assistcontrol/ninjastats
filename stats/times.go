package stats

import "math"

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
