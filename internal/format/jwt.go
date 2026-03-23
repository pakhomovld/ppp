package format

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/pakhomovld/pp/internal/color"
)

// JWTFormatter decodes and pretty-prints JWT header and payload.
type JWTFormatter struct{}

func (f *JWTFormatter) Format(w io.Writer, r io.Reader, theme *color.Theme) error {
	input, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	token := strings.TrimSpace(string(input))
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		_, writeErr := w.Write(input)
		return writeErr
	}

	headerLabel := "Header"
	payloadLabel := "Payload"
	sigLabel := "Signature"
	sigNote := "[not verified]"

	if theme != nil {
		headerLabel = theme.Sprint(color.Key, headerLabel)
		payloadLabel = theme.Sprint(color.Key, payloadLabel)
		sigLabel = theme.Sprint(color.Key, sigLabel)
		sigNote = theme.Sprint(color.Comment, sigNote)
	}

	// Header.
	fmt.Fprintf(w, "%s:\n", headerLabel)
	if err := printDecodedJSON(w, parts[0], theme); err != nil {
		fmt.Fprintf(w, "  (decode error: %v)\n", err)
	}

	fmt.Fprintln(w)

	// Payload.
	fmt.Fprintf(w, "%s:\n", payloadLabel)
	if err := printDecodedJSON(w, parts[1], theme); err != nil {
		fmt.Fprintf(w, "  (decode error: %v)\n", err)
	}

	fmt.Fprintln(w)

	// Signature (just show it exists).
	fmt.Fprintf(w, "%s: %s\n", sigLabel, sigNote)

	return nil
}

func printDecodedJSON(w io.Writer, b64 string, theme *color.Theme) error {
	decoded, err := base64.RawURLEncoding.DecodeString(b64)
	if err != nil {
		return err
	}

	var v any
	if err := json.Unmarshal(decoded, &v); err != nil {
		// Not JSON — just print raw decoded.
		fmt.Fprintf(w, "  %s\n", string(decoded))
		return nil
	}

	jsonFmt := &JSONFormatter{}
	// Indent the JSON output.
	return jsonFmt.Format(w, strings.NewReader(string(decoded)), theme)
}
