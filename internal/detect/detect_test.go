package detect

import "testing"

func TestDetect_JSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Format
	}{
		{"object", `{"key": "value"}`, JSON},
		{"array", `[1, 2, 3]`, JSON},
		{"nested", `{"a": {"b": [1, 2]}}`, JSON},
		{"empty object", `{}`, JSON},
		{"empty array", `[]`, JSON},
		{"with whitespace", `  { "key": 1 }  `, JSON},
		{"plain text", `hello world`, Plain},
		{"empty", ``, Plain},
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

func TestJSONDetector_Confidence(t *testing.T) {
	d := &JSONDetector{}

	tests := []struct {
		name     string
		input    string
		wantConf Confidence
	}{
		{"valid object", `{"a": 1}`, High},
		{"valid array", `[1, 2]`, High},
		{"truncated", `{"a": 1, "b":`, Medium},
		{"not json", `hello world`, None},
		{"empty", ``, None},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := d.Detect([]byte(tt.input))
			if r.Confidence != tt.wantConf {
				t.Errorf("confidence for %q = %v, want %v", tt.input, r.Confidence, tt.wantConf)
			}
		})
	}
}
