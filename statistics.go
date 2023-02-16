package main

import "math"

const (
	StdDevs = 3    // How many deviations from the mean to accept
	Zscore  = 1.96 // 95% confidence interval
)

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

// CIWidth returns the width of the confidence interval of
// the population mean of the times
func (ts Times) CIWidth() float64 {
	N := float64(ts.Count())
	stddev := ts.stddev()

	if N <= 0 {
		return 0
	}

	return Zscore * (stddev / math.Sqrt(N))
}

// stddev returns the population standard deviation of the times
func (ts Times) stddev() float64 {
	N := float64(ts.Count())
	variance := ts.variance()

	if N == 0 || variance < 0 {
		return 0
	}

	return math.Sqrt(variance / N)
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
