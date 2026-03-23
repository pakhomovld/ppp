package detect

import (
	"bytes"
	"regexp"
)

var (
	tomlSectionRe  = regexp.MustCompile(`(?m)^\s*\[[\w.\-]+\]\s*$`)
	tomlKeyValueRe = regexp.MustCompile(`(?m)^[\w.\-]+\s*=\s*.+`)
)

// TOMLDetector detects TOML documents.
type TOMLDetector struct{}

func (d *TOMLDetector) Detect(sample []byte) Result {
	trimmed := bytes.TrimSpace(sample)
	if len(trimmed) == 0 {
		return Result{Format: TOML, Confidence: None}
	}

	// Skip JSON and XML.
	if trimmed[0] == '{' {
		return Result{Format: TOML, Confidence: None}
	}
	if trimmed[0] == '[' && (len(trimmed) < 2 || trimmed[1] != '[') {
		// Bare `[` could be TOML section — let it fall through to regex check.
	}
	if trimmed[0] == '<' {
		return Result{Format: TOML, Confidence: None}
	}

	hasSections := tomlSectionRe.Match(trimmed)
	kvMatches := tomlKeyValueRe.FindAll(trimmed, -1)
	hasKeyValues := len(kvMatches) >= 1

	// TOML uses `=`, YAML uses `:`. Check that we have `=` assignments.
	if hasSections && hasKeyValues {
		return Result{Format: TOML, Confidence: High}
	}
	if hasKeyValues && len(kvMatches) >= 3 {
		return Result{Format: TOML, Confidence: Medium}
	}

	return Result{Format: TOML, Confidence: None}
}
