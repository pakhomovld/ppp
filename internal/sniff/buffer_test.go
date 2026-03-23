package sniff

import (
	"io"
	"strings"
	"testing"
)

func TestNewReader_SmallInput(t *testing.T) {
	input := "hello world"
	sr, err := NewReader(strings.NewReader(input), DefaultSize)
	if err != nil {
		t.Fatal(err)
	}

	sample := sr.Sample()
	if string(sample) != input {
		t.Errorf("Sample() = %q, want %q", sample, input)
	}

	all, err := io.ReadAll(sr.Reader())
	if err != nil {
		t.Fatal(err)
	}
	if string(all) != input {
		t.Errorf("Reader() produced %q, want %q", all, input)
	}
}

func TestNewReader_LargeInput(t *testing.T) {
	input := strings.Repeat("x", 20000)
	sr, err := NewReader(strings.NewReader(input), DefaultSize)
	if err != nil {
		t.Fatal(err)
	}

	sample := sr.Sample()
	if len(sample) != DefaultSize {
		t.Errorf("Sample() len = %d, want %d", len(sample), DefaultSize)
	}

	all, err := io.ReadAll(sr.Reader())
	if err != nil {
		t.Fatal(err)
	}
	if string(all) != input {
		t.Errorf("Reader() produced %d bytes, want %d", len(all), len(input))
	}
}

func TestNewReader_ExactBoundary(t *testing.T) {
	// Exactly DefaultSize bytes.
	input := strings.Repeat("a", DefaultSize)
	sr, err := NewReader(strings.NewReader(input), DefaultSize)
	if err != nil {
		t.Fatal(err)
	}

	sample := sr.Sample()
	if len(sample) != DefaultSize {
		t.Errorf("Sample() len = %d, want %d", len(sample), DefaultSize)
	}

	all, err := io.ReadAll(sr.Reader())
	if err != nil {
		t.Fatal(err)
	}
	if string(all) != input {
		t.Errorf("Reader() produced %d bytes, want %d", len(all), len(input))
	}
}

func TestNewReader_Empty(t *testing.T) {
	sr, err := NewReader(strings.NewReader(""), DefaultSize)
	if err != nil {
		t.Fatal(err)
	}

	if len(sr.Sample()) != 0 {
		t.Errorf("Sample() should be empty for empty input")
	}
}
