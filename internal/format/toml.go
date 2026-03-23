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
	tomlSectionLineRe = regexp.MustCompile(`^(\s*)(\[{1,2}[^\]]+\]{1,2})(.*)$`)
	tomlKVLineRe      = regexp.MustCompile(`^(\s*)([\w.\-]+)(\s*=\s*)(.+)$`)
	tomlCommentLineRe = regexp.MustCompile(`^(\s*)(#.*)$`)
)

// TOMLFormatter colorizes TOML output line by line.
type TOMLFormatter struct{}

func (f *TOMLFormatter) Format(w io.Writer, r io.Reader, theme *color.Theme) error {
	if theme == nil {
		_, err := io.Copy(w, r)
		return err
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		colored := colorizeTOMLLine(line, theme)
		if _, err := fmt.Fprintln(w, colored); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func colorizeTOMLLine(line string, theme *color.Theme) string {
	// Comment.
	if m := tomlCommentLineRe.FindStringSubmatch(line); m != nil {
		return m[1] + theme.Sprint(color.Comment, m[2])
	}

	// Section header [section] or [[array]].
	if m := tomlSectionLineRe.FindStringSubmatch(line); m != nil {
		result := m[1] + theme.Sprint(color.Key, m[2])
		if m[3] != "" {
			result += theme.Sprint(color.Comment, m[3])
		}
		return result
	}

	// Key = value.
	if m := tomlKVLineRe.FindStringSubmatch(line); m != nil {
		indent, key, sep, val := m[1], m[2], m[3], m[4]
		return indent + theme.Sprint(color.Key, key) +
			theme.Sprint(color.Colon, sep) +
			colorizeTOMLValue(val, theme)
	}

	return line
}

func colorizeTOMLValue(val string, theme *color.Theme) string {
	trimmed := strings.TrimSpace(val)

	switch {
	case trimmed == "true" || trimmed == "false":
		return theme.Sprint(color.Boolean, val)
	case isQuotedString(trimmed):
		return theme.Sprint(color.String, val)
	case isNumeric(trimmed):
		return theme.Sprint(color.Number, val)
	case strings.HasPrefix(trimmed, "["):
		return theme.Sprint(color.Bracket, val)
	case strings.HasPrefix(trimmed, "{"):
		return theme.Sprint(color.Bracket, val)
	default:
		return theme.Sprint(color.String, val)
	}
}
