package format

import (
	"bytes"
	"strings"
	"testing"
)

func TestXMLFormatter_NoColor(t *testing.T) {
	input := `<root><child id="1"><name>Alice</name></child></root>`

	f := &XMLFormatter{}
	var buf bytes.Buffer
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "  <child") {
		t.Error("expected indented child element")
	}
	if !strings.Contains(out, "    <name>") || !strings.Contains(out, "Alice") {
		t.Error("expected indented name element with content")
	}
}

func TestXMLFormatter_MalformedPassthrough(t *testing.T) {
	f := &XMLFormatter{}
	var buf bytes.Buffer
	input := "<root><unclosed"
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "<root>") && !strings.Contains(buf.String(), input) {
		t.Error("malformed XML should pass through or partially render")
	}
}

func TestXMLFormatter_SelfClosingTags(t *testing.T) {
	f := &XMLFormatter{}
	var buf bytes.Buffer
	input := `<root><item id="1"/><item id="2"/></root>`
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "item") {
		t.Error("self-closing tags should be present in output")
	}
}
