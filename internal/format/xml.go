package format

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/pakhomovld/pp/internal/color"
)

// XMLFormatter indents and colorizes XML/HTML.
type XMLFormatter struct{}

func (f *XMLFormatter) Format(w io.Writer, r io.Reader, theme *color.Theme) error {
	decoder := xml.NewDecoder(r)
	decoder.Strict = false
	decoder.AutoClose = xml.HTMLAutoClose

	depth := 0

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Malformed XML — read remaining and passthrough.
			return nil
		}

		indent := strings.Repeat("  ", depth)

		switch t := tok.(type) {
		case xml.StartElement:
			line := indent + formatStartElement(t, theme)
			fmt.Fprintln(w, line)
			depth++

		case xml.EndElement:
			depth--
			if depth < 0 {
				depth = 0
			}
			indent = strings.Repeat("  ", depth)
			tag := "</" + t.Name.Local + ">"
			if theme != nil {
				tag = theme.Sprint(color.Bracket, "</") +
					theme.Sprint(color.Key, t.Name.Local) +
					theme.Sprint(color.Bracket, ">")
			}
			fmt.Fprintln(w, indent+tag)

		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" {
				if theme != nil {
					text = theme.Sprint(color.String, text)
				}
				fmt.Fprintln(w, indent+text)
			}

		case xml.Comment:
			comment := "<!--" + string(t) + "-->"
			if theme != nil {
				comment = theme.Sprint(color.Comment, comment)
			}
			fmt.Fprintln(w, indent+comment)

		case xml.ProcInst:
			pi := "<?" + t.Target + " " + string(t.Inst) + "?>"
			if theme != nil {
				pi = theme.Sprint(color.Bracket, "<?") +
					theme.Sprint(color.Key, t.Target) + " " +
					theme.Sprint(color.String, string(t.Inst)) +
					theme.Sprint(color.Bracket, "?>")
			}
			fmt.Fprintln(w, indent+pi)

		case xml.Directive:
			dir := "<!" + string(t) + ">"
			if theme != nil {
				dir = theme.Sprint(color.Comment, dir)
			}
			fmt.Fprintln(w, indent+dir)
		}
	}

	return nil
}

func formatStartElement(el xml.StartElement, theme *color.Theme) string {
	var b strings.Builder

	if theme != nil {
		b.WriteString(theme.Sprint(color.Bracket, "<"))
		b.WriteString(theme.Sprint(color.Key, el.Name.Local))
	} else {
		b.WriteString("<")
		b.WriteString(el.Name.Local)
	}

	for _, attr := range el.Attr {
		b.WriteByte(' ')
		if theme != nil {
			b.WriteString(theme.Sprint(color.Number, attr.Name.Local))
			b.WriteString(theme.Sprint(color.Bracket, "="))
			b.WriteString(theme.Sprint(color.String, `"`+attr.Value+`"`))
		} else {
			b.WriteString(attr.Name.Local)
			b.WriteString(`="`)
			b.WriteString(attr.Value)
			b.WriteByte('"')
		}
	}

	if theme != nil {
		b.WriteString(theme.Sprint(color.Bracket, ">"))
	} else {
		b.WriteString(">")
	}

	return b.String()
}
