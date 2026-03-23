package format

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/pakhomovld/pp/internal/color"
)

// JSONFormatter pretty-prints JSON with indentation and color.
type JSONFormatter struct{}

func (f *JSONFormatter) Format(w io.Writer, r io.Reader, theme *color.Theme) error {
	input, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	var v any
	if err := json.Unmarshal(input, &v); err != nil {
		// Fall back to passthrough on invalid JSON.
		_, writeErr := w.Write(input)
		return writeErr
	}

	out := formatValue(v, 0, theme)
	_, err = fmt.Fprintln(w, out)
	return err
}

func formatValue(v any, depth int, theme *color.Theme) string {
	switch val := v.(type) {
	case nil:
		return theme.Sprint(color.Null, "null")
	case bool:
		if val {
			return theme.Sprint(color.Boolean, "true")
		}
		return theme.Sprint(color.Boolean, "false")
	case float64:
		return theme.Sprint(color.Number, formatNumber(val))
	case string:
		escaped, _ := json.Marshal(val)
		return theme.Sprint(color.String, string(escaped))
	case []any:
		return formatArray(val, depth, theme)
	case map[string]any:
		return formatObject(val, depth, theme)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func formatNumber(f float64) string {
	if f == float64(int64(f)) {
		return fmt.Sprintf("%d", int64(f))
	}
	return fmt.Sprintf("%g", f)
}

func formatArray(arr []any, depth int, theme *color.Theme) string {
	if len(arr) == 0 {
		return theme.Sprint(color.Bracket, "[]")
	}

	indent := strings.Repeat("  ", depth+1)
	closingIndent := strings.Repeat("  ", depth)

	var b strings.Builder
	b.WriteString(theme.Sprint(color.Bracket, "["))
	b.WriteByte('\n')

	for i, item := range arr {
		b.WriteString(indent)
		b.WriteString(formatValue(item, depth+1, theme))
		if i < len(arr)-1 {
			b.WriteString(theme.Sprint(color.Comma, ","))
		}
		b.WriteByte('\n')
	}

	b.WriteString(closingIndent)
	b.WriteString(theme.Sprint(color.Bracket, "]"))
	return b.String()
}

func formatObject(obj map[string]any, depth int, theme *color.Theme) string {
	if len(obj) == 0 {
		return theme.Sprint(color.Bracket, "{}")
	}

	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	indent := strings.Repeat("  ", depth+1)
	closingIndent := strings.Repeat("  ", depth)

	var b strings.Builder
	b.WriteString(theme.Sprint(color.Bracket, "{"))
	b.WriteByte('\n')

	for i, k := range keys {
		escaped, _ := json.Marshal(k)
		b.WriteString(indent)
		b.WriteString(theme.Sprint(color.Key, string(escaped)))
		b.WriteString(theme.Sprint(color.Colon, ": "))
		b.WriteString(formatValue(obj[k], depth+1, theme))
		if i < len(keys)-1 {
			b.WriteString(theme.Sprint(color.Comma, ","))
		}
		b.WriteByte('\n')
	}

	b.WriteString(closingIndent)
	b.WriteString(theme.Sprint(color.Bracket, "}"))
	return b.String()
}
