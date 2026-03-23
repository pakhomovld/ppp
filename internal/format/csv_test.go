package format

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pakhomovld/pp/internal/detect"
)

func TestCSVFormatter_NoColor(t *testing.T) {
	input := "name,age,city\nAlice,30,NYC\nBob,25,LA\n"

	f := &CSVFormatter{Dialect: detect.CSV}
	var buf bytes.Buffer
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	// Should contain header and data rows as aligned table.
	if !strings.Contains(out, "name") {
		t.Error("output should contain header 'name'")
	}
	if !strings.Contains(out, "Alice") {
		t.Error("output should contain data 'Alice'")
	}
	if !strings.Contains(out, "─") {
		t.Error("output should contain separator line")
	}
}

func TestCSVFormatter_TSV(t *testing.T) {
	input := "name\tage\nAlice\t30\nBob\t25\n"

	f := &CSVFormatter{Dialect: detect.TSV}
	var buf bytes.Buffer
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "Alice") {
		t.Error("output should contain 'Alice'")
	}
}
