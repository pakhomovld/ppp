package format

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/pakhomovld/pp/internal/color"
	"github.com/pakhomovld/pp/internal/detect"
)

// CSVFormatter renders CSV/TSV data as an aligned table.
type CSVFormatter struct {
	Dialect detect.Format
}

func (f *CSVFormatter) Format(w io.Writer, r io.Reader, theme *color.Theme) error {
	reader := csv.NewReader(r)
	if f.Dialect == detect.TSV {
		reader.Comma = '\t'
	}
	reader.FieldsPerRecord = -1 // Allow variable field count.
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		// Fall back: just copy raw input (can't re-read r though).
		return fmt.Errorf("csv parse error: %w", err)
	}

	if len(records) == 0 {
		return nil
	}

	// Compute column widths.
	widths := computeWidths(records)

	// Print header.
	if len(records) > 0 {
		printRow(w, records[0], widths, theme, true)
		printSeparator(w, widths, theme)
	}

	// Print data rows.
	for _, row := range records[1:] {
		printRow(w, row, widths, theme, false)
	}

	return nil
}

func computeWidths(records [][]string) []int {
	maxCols := 0
	for _, row := range records {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	widths := make([]int, maxCols)
	for _, row := range records {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Cap column width at 40 for readability.
	for i := range widths {
		if widths[i] > 40 {
			widths[i] = 40
		}
	}

	return widths
}

func printRow(w io.Writer, row []string, widths []int, theme *color.Theme, isHeader bool) {
	cells := make([]string, len(widths))
	for i := range widths {
		val := ""
		if i < len(row) {
			val = row[i]
			if len(val) > widths[i] {
				val = val[:widths[i]-1] + "…"
			}
		}
		padded := val + strings.Repeat(" ", widths[i]-len(val))
		if theme != nil && isHeader {
			padded = theme.Sprint(color.Key, padded)
		}
		cells[i] = padded
	}
	fmt.Fprintf(w, "  %s\n", strings.Join(cells, "  "))
}

func printSeparator(w io.Writer, widths []int, theme *color.Theme) {
	parts := make([]string, len(widths))
	for i, width := range widths {
		parts[i] = strings.Repeat("─", width)
	}
	sep := "  " + strings.Join(parts, "──")
	if theme != nil {
		sep = theme.Sprint(color.Bracket, sep)
	}
	fmt.Fprintln(w, sep)
}
