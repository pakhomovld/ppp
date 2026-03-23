package format

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/pakhomovld/pp/internal/color"
	"github.com/pakhomovld/pp/internal/detect"
)

// Base64Formatter decodes base64 and recursively detects/formats the decoded content.
type Base64Formatter struct{}

func (f *Base64Formatter) Format(w io.Writer, r io.Reader, theme *color.Theme) error {
	input, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	trimmed := bytes.TrimSpace(input)
	// Join multi-line base64.
	joined := strings.ReplaceAll(string(trimmed), "\n", "")
	joined = strings.ReplaceAll(joined, "\r", "")

	decoded, err := base64.StdEncoding.DecodeString(joined)
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(joined)
		if err != nil {
			// Can't decode — passthrough.
			_, writeErr := w.Write(input)
			return writeErr
		}
	}

	if !utf8.Valid(decoded) {
		// Binary content — passthrough original.
		_, writeErr := w.Write(input)
		return writeErr
	}

	// Detect the inner format.
	innerFormat := detect.Detect(decoded)

	label := fmt.Sprintf("[base64 decoded → %s]", innerFormat)
	if theme != nil {
		label = theme.Sprint(color.Comment, label)
	}
	fmt.Fprintln(w, label)

	formatter := ForFormat(innerFormat)
	return formatter.Format(w, strings.NewReader(string(decoded)), theme)
}
