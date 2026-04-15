package ingestion

import "strings"

// PreparedCode is the normalized source used by detector heuristics.
type PreparedCode struct {
	Raw         string
	Compact     string
	Lower       string
	LineCount   int
	NewlineRate float64
}

// Prepare normalizes code so detectors can share stable text features.
func Prepare(code string) PreparedCode {
	raw := strings.ReplaceAll(code, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")

	lines := strings.Count(raw, "\n") + 1
	if raw == "" {
		lines = 0
	}

	compact := strings.Join(strings.Fields(raw), " ")
	lower := strings.ToLower(raw)

	rate := 0.0
	if len(raw) > 0 {
		rate = float64(strings.Count(raw, "\n")) / float64(len(raw))
	}

	return PreparedCode{
		Raw:         raw,
		Compact:     compact,
		Lower:       lower,
		LineCount:   lines,
		NewlineRate: rate,
	}
}
