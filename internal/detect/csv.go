package detect

import (
	"bytes"
	"encoding/csv"
	"io"
	"strings"
)

// CSVDetector detects CSV and TSV data by checking for consistent
// field counts across the first several lines.
type CSVDetector struct{}

func (d *CSVDetector) Detect(sample []byte) Result {
	trimmed := bytes.TrimSpace(sample)
	if len(trimmed) == 0 {
		return Result{Format: CSV, Confidence: None}
	}

	// Skip if it looks like JSON or XML.
	if trimmed[0] == '{' || trimmed[0] == '[' || trimmed[0] == '<' {
		return Result{Format: CSV, Confidence: None}
	}

	// Try comma first (RFC 4180 aware), then tab.
	if f, c := checkDelimiter(trimmed, ','); c > None {
		return Result{Format: f, Confidence: c}
	}
	if f, c := checkDelimiter(trimmed, '\t'); c > None {
		return Result{Format: f, Confidence: c}
	}

	return Result{Format: CSV, Confidence: None}
}

func checkDelimiter(sample []byte, delim byte) (Format, Confidence) {
	format := CSV
	if delim == '\t' {
		format = TSV
	}

	// Use encoding/csv for proper RFC 4180 field counting (handles quoted fields).
	reader := csv.NewReader(bytes.NewReader(sample))
	reader.Comma = rune(delim)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	// Read up to 10 records.
	maxRecords := 10
	var counts []int
	for len(counts) < maxRecords {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Parse error — not valid CSV with this delimiter.
			return format, None
		}
		line := strings.TrimSpace(strings.Join(record, string(delim)))
		if line == "" {
			continue
		}
		counts = append(counts, len(record))
	}

	if len(counts) < 2 {
		return format, None
	}

	// All lines must have the same field count, and it must be > 1.
	expected := counts[0]
	if expected <= 1 {
		return format, None
	}

	consistent := 0
	for _, c := range counts[1:] {
		if c == expected {
			consistent++
		}
	}

	total := len(counts) - 1
	if consistent == total && total >= 4 {
		return format, High
	}
	if consistent == total && total >= 1 {
		return format, Medium
	}

	return format, None
}
