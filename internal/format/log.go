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
	logTimestampRe = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}[.\d]*)(\S*)`)
	logLevelRe     = regexp.MustCompile(`(?i)\b(INFO|ERROR|WARN|WARNING|DEBUG|TRACE|FATAL)\b`)
	logBracketedRe = regexp.MustCompile(`(?i)\[(INFO|ERROR|WARN|WARNING|DEBUG|TRACE|FATAL)\]`)
)

// LogFormatter colorizes log lines, processing them one at a time (streaming).
type LogFormatter struct{}

func (f *LogFormatter) Format(w io.Writer, r io.Reader, theme *color.Theme) error {
	if theme == nil {
		_, err := io.Copy(w, r)
		return err
	}

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		colored := colorizeLogLine(line, theme)
		if _, err := fmt.Fprintln(w, colored); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func colorizeLogLine(line string, theme *color.Theme) string {
	result := line

	// Colorize timestamp.
	result = logTimestampRe.ReplaceAllStringFunc(result, func(match string) string {
		return theme.Sprint(color.Comment, match)
	})

	// Colorize bracketed levels: [ERROR], [INFO], etc.
	result = logBracketedRe.ReplaceAllStringFunc(result, func(match string) string {
		return colorizeLevel(match, theme)
	})

	// Colorize bare levels: ERROR, INFO, etc. (only if not already colored).
	if !logBracketedRe.MatchString(line) {
		result = logLevelRe.ReplaceAllStringFunc(result, func(match string) string {
			return colorizeLevel(match, theme)
		})
	}

	return result
}

func colorizeLevel(level string, theme *color.Theme) string {
	upper := strings.ToUpper(level)
	switch {
	case strings.Contains(upper, "ERROR"), strings.Contains(upper, "FATAL"):
		return theme.Sprint(color.Null, level) // Red.
	case strings.Contains(upper, "WARN"):
		return theme.Sprint(color.Number, level) // Yellow.
	case strings.Contains(upper, "INFO"):
		return theme.Sprint(color.Boolean, level) // Magenta (stands out).
	case strings.Contains(upper, "DEBUG"), strings.Contains(upper, "TRACE"):
		return theme.Sprint(color.Comment, level) // Dim.
	default:
		return level
	}
}
