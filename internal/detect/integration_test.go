package detect

import "testing"

func TestDetect_FormatPriority(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Format
	}{
		// JSON should win over YAML (valid JSON is valid YAML).
		{"json over yaml", `{"key": "value"}`, JSON},
		{"json array", `[1, 2, 3]`, JSON},

		// YAML detection.
		{"yaml doc separator", "---\nname: test\nage: 30", YAML},
		{"yaml key-value", "name: test\nage: 30\ncity: NYC", YAML},

		// CSV detection.
		{"csv", "a,b,c\n1,2,3\n4,5,6\n7,8,9\n10,11,12", CSV},
		{"tsv", "a\tb\tc\n1\t2\t3\n4\t5\t6\n7\t8\t9\n10\t11\t12", TSV},

		// XML/HTML detection.
		{"xml", `<?xml version="1.0"?><root/>`, XML},
		{"html", `<!DOCTYPE html><html></html>`, HTML},
		{"xml tag", `<root><child>text</child></root>`, XML},

		// TOML detection.
		{"toml", "[server]\nhost = \"localhost\"\nport = 8080", TOML},

		// Log detection.
		{"logs", "2024-01-15T10:30:00Z INFO Start\n2024-01-15T10:30:01Z ERROR Fail\n2024-01-15T10:30:02Z WARN Retry", LogLine},

		// JWT detection (must beat JSON since JWT starts with eyJ which is valid base64 of JSON).
		{"jwt", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U", JWT},

		// URL-encoded detection.
		{"url-encoded", "name=Alice&age=30&city=NYC", URLEncode},

		// Fallbacks.
		{"plain text", "just some text", Plain},
		{"empty", "", Plain},
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
