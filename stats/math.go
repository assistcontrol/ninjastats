package stats

import "math"

const (
	stdDevs = 3    // How many deviations from the mean to accept
	zScore  = 1.96 // 95% confidence interval
)

// count returns the total number of times (including outliers)
func (ts times) count() int {
	return len(ts.times)
}

// mean returns the mean time, or 0 for an empty list
func (ts times) mean() float64 {
	N := ts.count()
	if N == 0 {
		return 0
	}

	return ts.sum() / float64(ts.count())
}

// ciWidth returns the width of the confidence interval of
// the population mean of the times
func (ts times) ciWidth() float64 {
	N := float64(ts.count())
	stddev := ts.stddev()

	if N <= 0 {
		return 0
	}

	return zScore * (stddev / math.Sqrt(N))
}

// reductionLimit returns the limit that reduce() should use
func (ts times) reductionLimit() float64 {
	return stdDevs * ts.stddev()
}

// stddev returns the population standard deviation of the times
func (ts times) stddev() float64 {
	N := float64(ts.count())
	variance := ts.variance()

	if N == 0 || variance < 0 {
		return 0
	}

	return math.Sqrt(variance / N)
}

// sum returns the sum of all times
func (ts times) sum() (sum float64) {
	for _, t := range ts.times {
		sum += t
	}

	return
}

// variance returns the population variance of the times
func (ts times) variance() (v float64) {
	mean := ts.mean()
	for _, t := range ts.times {
		v += (t - mean) * (t - mean)
	}

	return
}
