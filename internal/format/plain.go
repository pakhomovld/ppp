package format

import (
	"io"

	"github.com/pakhomovld/pp/internal/color"
)

// PlainFormatter passes input through without modification.
type PlainFormatter struct{}

func (f *PlainFormatter) Format(w io.Writer, r io.Reader, _ *color.Theme) error {
	_, err := io.Copy(w, r)
	return err
}
