package format

import "github.com/pakhomovld/pp/internal/detect"

// ForFormat returns the appropriate formatter for the given format.
func ForFormat(f detect.Format) Formatter {
	switch f {
	case detect.JSON:
		return &JSONFormatter{}
	case detect.YAML:
		return &YAMLFormatter{}
	case detect.CSV:
		return &CSVFormatter{Dialect: detect.CSV}
	case detect.TSV:
		return &CSVFormatter{Dialect: detect.TSV}
	case detect.TOML:
		return &TOMLFormatter{}
	case detect.XML, detect.HTML:
		return &XMLFormatter{}
	case detect.LogLine:
		return &LogFormatter{}
	case detect.JWT:
		return &JWTFormatter{}
	case detect.Base64:
		return &Base64Formatter{}
	case detect.URLEncode:
		return &URLFormatter{}
	default:
		return &PlainFormatter{}
	}
}
