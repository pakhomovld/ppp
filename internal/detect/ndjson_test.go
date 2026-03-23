package detect

import "testing"

func TestNDJSONDetector(t *testing.T) {
	d := &NDJSONDetector{}

	tests := []struct {
		name     string
		input    string
		wantConf Confidence
	}{
		{"two valid objects", "{\"a\":1}\n{\"b\":2}\n", High},
		{"three valid objects", "{\"a\":1}\n{\"b\":2}\n{\"c\":3}", High},
		{"valid with empty lines", "{\"a\":1}\n\n{\"b\":2}\n", High},
		{"mixed valid and truncated", "{\"a\":1}\n{\"b\":2,\"c\":", Medium},
		{"single line", "{\"a\":1}", None},
		{"empty", "", None},
		{"lines starting with bracket", "[1,2]\n[3,4]", None},
		{"plain text lines", "hello\nworld", None},
		{"mixed json and text", "{\"a\":1}\nhello", None},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := d.Detect([]byte(tt.input))
			if r.Confidence != tt.wantConf {
				t.Errorf("NDJSONDetector(%q).Confidence = %v, want %v", tt.input, r.Confidence, tt.wantConf)
			}
			if r.Confidence != None && r.Format != NDJSON {
				t.Errorf("NDJSONDetector(%q).Format = %v, want %v", tt.input, r.Format, NDJSON)
			}
		})
	}
}
