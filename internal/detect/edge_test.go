package detect

import (
	"strings"
	"testing"
)

func TestEdge_YAMLvsProse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Format
	}{
		{"letter greeting", "Dear John: hope you are well", Plain},
		{"short note", "Note: important", Plain},
		{"single word colon", "hello: world", Plain},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Detect([]byte(tt.input)).Format
			if got != tt.want {
				t.Errorf("Detect(%q).Format = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestEdge_BinaryGarbage(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  Format
	}{
		{"null bytes", []byte{0x00, 0x01, 0x02}, Plain},
		{"high bytes", []byte{0xFF, 0xFE, 0xFD, 0xFC}, Plain},
		{"mixed garbage", []byte{0x80, 0x90, 0xA0, 0xB0, 0xC0}, Plain},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Detect(tt.input).Format
			if got != tt.want {
				t.Errorf("Detect(%v).Format = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestEdge_TruncatedJWT(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Format
	}{
		{"two segments only", "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxIn0", Plain},
		{"empty segment", "eyJhbGciOiJIUzI1NiJ9..sig", Plain},
		{"four segments", "a.b.c.d", Plain},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Detect([]byte(tt.input)).Format
			if got == JWT {
				t.Errorf("Detect(%q).Format = %q, should not be JWT", tt.input, got)
			}
		})
	}
}

func TestEdge_TruncatedBase64(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"too short", "YWJj"},
		{"bad padding", "YWJj====extra"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Detect([]byte(tt.input)).Format
			if got == Base64 {
				t.Errorf("Detect(%q).Format = base64, should not be base64 for truncated/short input", tt.input)
			}
		})
	}
}

func TestEdge_BrokenXML(t *testing.T) {
	input := "<root><unclosed"
	r := Detect([]byte(input))
	if r.Format != XML {
		t.Errorf("Detect(%q).Format = %q, want xml", input, r.Format)
	}
	if r.Confidence > Medium {
		t.Errorf("Detect(%q).Confidence = %v, want <= Medium", input, r.Confidence)
	}
}

func TestEdge_CSVQuotedFields(t *testing.T) {
	// Fields with commas inside quotes — CSV detection may treat this as more columns.
	input := `name,city,note
"Alice","New York","hello, world"
"Bob","LA","ok"
"Carol","SF","fine"
"Dave","CHI","great"`
	r := Detect([]byte(input))
	// Document: naive comma counting may miscount quoted fields.
	// At minimum, it should not panic.
	if r.Format == "" {
		t.Error("Detect should return a format for quoted CSV")
	}
}

func TestEdge_WhitespaceOnly(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"spaces", "     "},
		{"newlines", "\n\n\n"},
		{"tabs", "\t\t\t"},
		{"mixed", "  \n\t  \n  "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Detect([]byte(tt.input)).Format
			if got != Plain {
				t.Errorf("Detect(whitespace).Format = %q, want plain", got)
			}
		})
	}
}

func TestEdge_TOMLSingleBracket(t *testing.T) {
	// Should not panic on single bracket (was a bug).
	input := "["
	r := Detect([]byte(input))
	_ = r // Must not panic.
}

func TestEdge_SniffBufferBoundary(t *testing.T) {
	// JSON at exactly 8192 bytes.
	padding := strings.Repeat(" ", 8192-len(`{"a":1}`))
	input8192 := `{"a":1}` + padding
	if len(input8192) != 8192 {
		t.Fatalf("expected 8192 bytes, got %d", len(input8192))
	}
	got := Detect([]byte(input8192)).Format
	if got != JSON {
		t.Errorf("Detect(8192-byte JSON).Format = %q, want json", got)
	}

	// JSON at 8193 bytes (one past boundary, but Detect gets the full sample here).
	input8193 := input8192 + " "
	got = Detect([]byte(input8193)).Format
	if got != JSON {
		t.Errorf("Detect(8193-byte JSON).Format = %q, want json", got)
	}
}
