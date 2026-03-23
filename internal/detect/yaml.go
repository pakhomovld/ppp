package detect

import (
	"bytes"
	"regexp"
)

var yamlKeyPattern = regexp.MustCompile(`(?m)^[a-zA-Z_][a-zA-Z0-9_.\-]*\s*:`)

// YAMLDetector detects YAML documents.
type YAMLDetector struct{}

func (d *YAMLDetector) Detect(sample []byte) Result {
	trimmed := bytes.TrimSpace(sample)
	if len(trimmed) == 0 {
		return Result{Format: YAML, Confidence: None}
	}

	// If it starts with '<', it's likely XML, not YAML.
	if trimmed[0] == '<' {
		return Result{Format: YAML, Confidence: None}
	}

	// If it's valid JSON, let the JSON detector win.
	if trimmed[0] == '{' || trimmed[0] == '[' {
		return Result{Format: YAML, Confidence: None}
	}

	// YAML document separator is a strong signal.
	if bytes.HasPrefix(trimmed, []byte("---")) {
		return Result{Format: YAML, Confidence: High}
	}

	// Check for key: value patterns across multiple lines.
	lines := bytes.Split(trimmed, []byte("\n"))
	matches := 0
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		if yamlKeyPattern.Match(line) {
			matches++
		}
	}

	if matches >= 3 {
		return Result{Format: YAML, Confidence: Medium}
	}
	if matches == 2 {
		return Result{Format: YAML, Confidence: Low}
	}

	return Result{Format: YAML, Confidence: None}
}
