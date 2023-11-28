package format

import (
	"fmt"

	"github.com/assistcontrol/ninjastats/stats"
)

const (
	HeaderCI    = "msec/page ± 95% CI:"
	HeaderCount = "count (outliers):"
	HeaderTimes = "msec/page (pages/sec):"
)

// Type Formatter is a function that formats a single verb's Report
// as a string (< 16 chars).
type Formatter func(r *stats.Report) string

func Time(r *stats.Report) string {
	// If Mean == 0, rate is nonsense
	if r.Mean == 0 {
		return ""
	}

	return fmt.Sprintf("%.2f (%d)", r.Mean, r.Rate)
}

func Count(r *stats.Report) string {
	// Reduce table clutter
	if r.Count == 0 {
		return ""
	}

	return fmt.Sprintf("%d (%5d)", r.Count, r.Outliers)
}

func CI(r *stats.Report) string {
	if r.Mean == 0 {
		return ""
	}

	return fmt.Sprintf("%.2f ± %.2f", r.Mean, r.CIWidth)
}
