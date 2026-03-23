package format

import (
	"fmt"
	"io"
	"net/url"
	"sort"
	"strings"

	"github.com/pakhomovld/pp/internal/color"
)

// URLFormatter decodes URL-encoded query strings and displays them as a table.
type URLFormatter struct{}

func (f *URLFormatter) Format(w io.Writer, r io.Reader, theme *color.Theme) error {
	input, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	query := strings.TrimSpace(string(input))
	values, err := url.ParseQuery(query)
	if err != nil {
		_, writeErr := w.Write(input)
		return writeErr
	}

	// Sort keys for consistent output.
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Find max key length for alignment.
	maxKeyLen := 0
	for _, k := range keys {
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}

	for _, k := range keys {
		for _, v := range values[k] {
			key := k + strings.Repeat(" ", maxKeyLen-len(k))
			if theme != nil {
				key = theme.Sprint(color.Key, key)
				v = theme.Sprint(color.String, v)
			}
			sep := " = "
			if theme != nil {
				sep = theme.Sprint(color.Colon, sep)
			}
			fmt.Fprintf(w, "  %s%s%s\n", key, sep, v)
		}
	}

	return nil
}
