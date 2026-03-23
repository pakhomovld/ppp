package format

import (
	"bytes"
	"strings"
	"testing"
)

func TestJWTFormatter_NoColor(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	f := &JWTFormatter{}
	var buf bytes.Buffer
	err := f.Format(&buf, strings.NewReader(token), nil)
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "Header:") {
		t.Error("output should contain 'Header:'")
	}
	if !strings.Contains(out, `"alg"`) {
		t.Error("output should contain decoded 'alg' field")
	}
	if !strings.Contains(out, "Payload:") {
		t.Error("output should contain 'Payload:'")
	}
	if !strings.Contains(out, "John Doe") {
		t.Error("output should contain decoded 'name' value")
	}
	if !strings.Contains(out, "Signature:") {
		t.Error("output should contain 'Signature:'")
	}
}

func TestJWTFormatter_MalformedPassthrough(t *testing.T) {
	f := &JWTFormatter{}
	var buf bytes.Buffer
	input := "not.a.jwt"
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}
	// Malformed JWT should still produce some output (decode error or passthrough).
	if buf.Len() == 0 {
		t.Error("malformed JWT should produce output")
	}
}

func TestJWTFormatter_TwoSegments(t *testing.T) {
	f := &JWTFormatter{}
	var buf bytes.Buffer
	input := "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxIn0"
	err := f.Format(&buf, strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}
	// Not a valid JWT (2 parts) — should passthrough.
	if buf.String() != input {
		t.Errorf("2-segment JWT should pass through, got %q", buf.String())
	}
}
