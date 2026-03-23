package format

import (
	"bytes"
	"strings"
	"testing"
)

func TestJSONFormatter_NoColor(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			"simple object",
			`{"b":2,"a":1}`,
			"{\n  \"a\": 1,\n  \"b\": 2\n}\n",
		},
		{
			"array",
			`[1,2,3]`,
			"[\n  1,\n  2,\n  3\n]\n",
		},
		{
			"empty object",
			`{}`,
			"{}\n",
		},
		{
			"empty array",
			`[]`,
			"[]\n",
		},
		{
			"nested",
			`{"a":{"b":1}}`,
			"{\n  \"a\": {\n    \"b\": 1\n  }\n}\n",
		},
	}

	f := &JSONFormatter{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := f.Format(&buf, strings.NewReader(tt.input), nil)
			if err != nil {
				t.Fatal(err)
			}
			got := buf.String()
			if got != tt.want {
				t.Errorf("got:\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}

func TestJSONFormatter_DeepNesting(t *testing.T) {
	f := &JSONFormatter{}
	var buf bytes.Buffer
	input := `{"a":{"b":{"c":{"d":"deep"}}}}`
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"deep"`) {
		t.Error("deep nested value should be present")
	}
}

func TestJSONFormatter_Unicode(t *testing.T) {
	f := &JSONFormatter{}
	var buf bytes.Buffer
	input := `{"emoji":"🎉","japanese":"日本語"}`
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "🎉") {
		t.Error("emoji should be preserved")
	}
	if !strings.Contains(out, "日本語") {
		t.Error("unicode text should be preserved")
	}
}

func TestJSONFormatter_NullValues(t *testing.T) {
	f := &JSONFormatter{}
	var buf bytes.Buffer
	input := `{"a":null,"b":[null,1,null]}`
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "null") {
		t.Error("null values should be preserved")
	}
}

func TestJSONFormatter_InvalidJSON(t *testing.T) {
	f := &JSONFormatter{}
	var buf bytes.Buffer
	input := "not json at all"
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != input {
		t.Errorf("invalid JSON should pass through unchanged, got %q", buf.String())
	}
}
