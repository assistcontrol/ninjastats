package stats

import "math"

// Times
type times struct {
	times []float64
}

// NewTimes creates the Times struct. Required for initialization of
// the underlying .times slice.
func newTimes() *times {
	return &times{
		times: make([]float64, 0),
	}
}

// Add appends a time tm to the collection of times
func (ts *times) add(tm float64) {
	ts.times = append(ts.times, tm)
}

// Reduce produces a new Times struct from the parent data with
// outliers removed
func (ts times) reduce() *times {
	mean := ts.mean()
	limit := ts.reductionLimit()

	reduced := &times{}
	for _, t := range ts.times {
		if math.Abs(t-mean) <= limit {
			reduced.add(t)
		}
	}

	return reduced
}
