package format

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/pakhomovld/pp/internal/color"
)

var (
	yamlKeyRe     = regexp.MustCompile(`^(\s*)([\w.\-/]+)(\s*:\s*)(.*)$`)
	yamlCommentRe = regexp.MustCompile(`^(\s*)(#.*)$`)
	yamlDashRe    = regexp.MustCompile(`^(\s*)(- )(.*)$`)
)

// YAMLFormatter colorizes YAML output line by line.
// YAML is already human-readable, so we just add color — no restructuring.
type YAMLFormatter struct{}

func (f *YAMLFormatter) Format(w io.Writer, r io.Reader, theme *color.Theme) error {
	if theme == nil {
		_, err := io.Copy(w, r)
		return err
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		colored := colorizeYAMLLine(line, theme)
		if _, err := fmt.Fprintln(w, colored); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func colorizeYAMLLine(line string, theme *color.Theme) string {
	// Document separator.
	if strings.TrimSpace(line) == "---" || strings.TrimSpace(line) == "..." {
		return theme.Sprint(color.Bracket, line)
	}

	// Comment.
	if m := yamlCommentRe.FindStringSubmatch(line); m != nil {
		return m[1] + theme.Sprint(color.Comment, m[2])
	}

	// Key: value.
	if m := yamlKeyRe.FindStringSubmatch(line); m != nil {
		indent, key, sep, val := m[1], m[2], m[3], m[4]
		return indent + theme.Sprint(color.Key, key) + theme.Sprint(color.Colon, sep) + colorizeYAMLValue(val, theme)
	}

	// List item.
	if m := yamlDashRe.FindStringSubmatch(line); m != nil {
		indent, dash, val := m[1], m[2], m[3]
		return indent + theme.Sprint(color.Bracket, dash) + colorizeYAMLValue(val, theme)
	}

	return line
}

func colorizeYAMLValue(val string, theme *color.Theme) string {
	trimmed := strings.TrimSpace(val)

	switch {
	case trimmed == "true" || trimmed == "false" ||
		trimmed == "yes" || trimmed == "no" ||
		trimmed == "on" || trimmed == "off":
		return theme.Sprint(color.Boolean, val)
	case trimmed == "null" || trimmed == "~" || trimmed == "":
		return theme.Sprint(color.Null, val)
	case isQuotedString(trimmed):
		return theme.Sprint(color.String, val)
	case isNumeric(trimmed):
		return theme.Sprint(color.Number, val)
	default:
		return theme.Sprint(color.String, val)
	}
}

func isQuotedString(s string) bool {
	return (strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) ||
		(strings.HasPrefix(s, `'`) && strings.HasSuffix(s, `'`))
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	dotCount := 0
	for i, c := range s {
		if c == '-' && i == 0 {
			continue
		}
		if c == '.' {
			dotCount++
			if dotCount > 1 {
				return false
			}
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
