package format

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"

	"github.com/pakhomovld/ppp/internal/color"
)

var sqlKeywords = map[string]bool{
	"SELECT": true, "FROM": true, "WHERE": true, "JOIN": true,
	"LEFT": true, "RIGHT": true, "INNER": true, "OUTER": true,
	"CROSS": true, "FULL": true, "ON": true, "AND": true,
	"OR": true, "NOT": true, "IN": true, "EXISTS": true,
	"BETWEEN": true, "LIKE": true, "IS": true, "NULL": true,
	"AS": true, "CASE": true, "WHEN": true, "THEN": true,
	"ELSE": true, "END": true, "INSERT": true, "INTO": true,
	"VALUES": true, "UPDATE": true, "SET": true, "DELETE": true,
	"CREATE": true, "TABLE": true, "INDEX": true, "VIEW": true,
	"ALTER": true, "DROP": true, "ADD": true, "COLUMN": true,
	"PRIMARY": true, "KEY": true, "FOREIGN": true, "REFERENCES": true,
	"DEFAULT": true, "CONSTRAINT": true, "GROUP": true, "BY": true,
	"ORDER": true, "ASC": true, "DESC": true, "HAVING": true,
	"LIMIT": true, "OFFSET": true, "UNION": true, "ALL": true,
	"DISTINCT": true, "EXPLAIN": true, "WITH": true, "BEGIN": true,
	"COMMIT": true, "ROLLBACK": true, "GRANT": true, "REVOKE": true,
	"COUNT": true, "SUM": true, "AVG": true, "MIN": true, "MAX": true,
	"DATABASE": true, "SCHEMA": true, "INTERSECT": true, "EXCEPT": true,
}

// Major clause keywords that get a newline before them.
var sqlNewlineKeywords = map[string]bool{
	"FROM": true, "WHERE": true, "JOIN": true, "LEFT": true,
	"RIGHT": true, "INNER": true, "OUTER": true, "CROSS": true,
	"FULL": true, "GROUP": true, "ORDER": true, "HAVING": true,
	"LIMIT": true, "UNION": true, "SET": true, "VALUES": true,
	"INTERSECT": true, "EXCEPT": true, "INTO": true, "ON": true,
	"AND": true, "OR": true,
}

var sqlCommentLineRe = regexp.MustCompile(`^\s*--`)

// SQLFormatter formats SQL queries with keyword uppercasing and indentation.
type SQLFormatter struct{}

func (f *SQLFormatter) Format(w io.Writer, r io.Reader, theme *color.Theme) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	input := strings.TrimSpace(string(data))
	tokens := tokenizeSQL(input)
	formatted := formatSQLTokens(tokens, theme)

	_, err = fmt.Fprintln(w, formatted)
	return err
}

type sqlTokenType int

const (
	sqlWord sqlTokenType = iota
	sqlString
	sqlNumber
	sqlPunctuation
	sqlComment
	sqlWhitespace
)

type sqlToken struct {
	typ   sqlTokenType
	value string
}

func tokenizeSQL(input string) []sqlToken {
	var tokens []sqlToken
	i := 0

	for i < len(input) {
		ch := input[i]

		// Line comment (--).
		if ch == '-' && i+1 < len(input) && input[i+1] == '-' {
			end := strings.IndexByte(input[i:], '\n')
			if end == -1 {
				tokens = append(tokens, sqlToken{sqlComment, input[i:]})
				break
			}
			tokens = append(tokens, sqlToken{sqlComment, input[i : i+end]})
			i += end
			continue
		}

		// Block comment (/* */).
		if ch == '/' && i+1 < len(input) && input[i+1] == '*' {
			end := strings.Index(input[i+2:], "*/")
			if end == -1 {
				tokens = append(tokens, sqlToken{sqlComment, input[i:]})
				break
			}
			tokens = append(tokens, sqlToken{sqlComment, input[i : i+2+end+2]})
			i += 2 + end + 2
			continue
		}

		// Single-quoted string.
		if ch == '\'' {
			end := i + 1
			for end < len(input) {
				if input[end] == '\'' {
					if end+1 < len(input) && input[end+1] == '\'' {
						end += 2 // escaped quote
						continue
					}
					end++
					break
				}
				end++
			}
			tokens = append(tokens, sqlToken{sqlString, input[i:end]})
			i = end
			continue
		}

		// Whitespace.
		if unicode.IsSpace(rune(ch)) {
			end := i + 1
			for end < len(input) && unicode.IsSpace(rune(input[end])) {
				end++
			}
			tokens = append(tokens, sqlToken{sqlWhitespace, input[i:end]})
			i = end
			continue
		}

		// Number.
		if ch >= '0' && ch <= '9' {
			end := i + 1
			for end < len(input) && (input[end] >= '0' && input[end] <= '9' || input[end] == '.') {
				end++
			}
			tokens = append(tokens, sqlToken{sqlNumber, input[i:end]})
			i = end
			continue
		}

		// Word (identifier or keyword).
		if isWordChar(ch) {
			end := i + 1
			for end < len(input) && isWordChar(input[end]) {
				end++
			}
			tokens = append(tokens, sqlToken{sqlWord, input[i:end]})
			i = end
			continue
		}

		// Punctuation / operators.
		tokens = append(tokens, sqlToken{sqlPunctuation, string(ch)})
		i++
	}

	return tokens
}

func isWordChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '.'
}

func formatSQLTokens(tokens []sqlToken, theme *color.Theme) string {
	var b strings.Builder
	indent := 1
	afterNewline := false
	isFirstToken := true

	for i, tok := range tokens {
		switch tok.typ {
		case sqlWhitespace:
			// Skip original whitespace; we control spacing.
			continue

		case sqlComment:
			if !isFirstToken && !afterNewline {
				b.WriteString(" ")
			}
			if theme != nil {
				b.WriteString(theme.Sprint(color.Comment, tok.value))
			} else {
				b.WriteString(tok.value)
			}
			b.WriteString("\n")
			writeIndent(&b, indent)
			afterNewline = true

		case sqlString:
			if !isFirstToken && !afterNewline {
				b.WriteString(" ")
			}
			if theme != nil {
				b.WriteString(theme.Sprint(color.String, tok.value))
			} else {
				b.WriteString(tok.value)
			}
			afterNewline = false

		case sqlNumber:
			if !isFirstToken && !afterNewline {
				b.WriteString(" ")
			}
			if theme != nil {
				b.WriteString(theme.Sprint(color.Number, tok.value))
			} else {
				b.WriteString(tok.value)
			}
			afterNewline = false

		case sqlPunctuation:
			switch tok.value {
			case "(":
				if !isFirstToken && !afterNewline {
					b.WriteString(" ")
				}
				if theme != nil {
					b.WriteString(theme.Sprint(color.Colon, tok.value))
				} else {
					b.WriteString(tok.value)
				}
				indent++
				afterNewline = false
			case ")":
				indent--
				if indent < 0 {
					indent = 0
				}
				if theme != nil {
					b.WriteString(theme.Sprint(color.Colon, tok.value))
				} else {
					b.WriteString(tok.value)
				}
				afterNewline = false
			case ",":
				if theme != nil {
					b.WriteString(theme.Sprint(color.Comma, tok.value))
				} else {
					b.WriteString(tok.value)
				}
				// After comma in SELECT list, add newline.
				if isInSelectList(tokens, i) {
					b.WriteString("\n")
					writeIndent(&b, indent)
					afterNewline = true
				}
			case ";":
				if theme != nil {
					b.WriteString(theme.Sprint(color.Colon, tok.value))
				} else {
					b.WriteString(tok.value)
				}
				b.WriteString("\n")
				afterNewline = true
			case "*", "=", "<", ">", "!", "+", "-", "/":
				if !isFirstToken && !afterNewline {
					b.WriteString(" ")
				}
				if theme != nil {
					b.WriteString(theme.Sprint(color.Colon, tok.value))
				} else {
					b.WriteString(tok.value)
				}
				afterNewline = false
			default:
				b.WriteString(tok.value)
				afterNewline = false
			}

		case sqlWord:
			upper := strings.ToUpper(tok.value)

			// Check if this keyword should trigger a newline.
			if !isFirstToken && sqlNewlineKeywords[upper] {
				// Don't newline for INTO when right after INSERT.
				if upper == "INTO" && isPrecededByKeyword(tokens, i, "INSERT") {
					// no newline
				} else {
					b.WriteString("\n")
					writeIndent(&b, indent-1)
					afterNewline = true
				}
			}

			if !isFirstToken && !afterNewline {
				b.WriteString(" ")
			}

			display := tok.value
			if sqlKeywords[upper] {
				display = upper
			}

			if theme != nil && sqlKeywords[upper] {
				b.WriteString(theme.Sprint(color.Key, display))
			} else {
				b.WriteString(display)
			}
			afterNewline = false
		}

		if tok.typ != sqlWhitespace {
			isFirstToken = false
		}
	}

	return b.String()
}

func writeIndent(b *strings.Builder, level int) {
	for range level {
		b.WriteString("  ")
	}
}

// isInSelectList checks if the comma at index i is within a SELECT column list
// (between SELECT and FROM).
func isInSelectList(tokens []sqlToken, commaIdx int) bool {
	// Look backwards for SELECT without hitting FROM.
	depth := 0
	for j := commaIdx - 1; j >= 0; j-- {
		if tokens[j].typ == sqlPunctuation {
			if tokens[j].value == ")" {
				depth++
			} else if tokens[j].value == "(" {
				depth--
			}
		}
		if depth > 0 {
			continue
		}
		if tokens[j].typ == sqlWord {
			upper := strings.ToUpper(tokens[j].value)
			if upper == "SELECT" {
				return true
			}
			if upper == "FROM" {
				return false
			}
		}
	}
	return false
}

// isPrecededByKeyword checks if the word token at index i is preceded by the given keyword.
func isPrecededByKeyword(tokens []sqlToken, i int, keyword string) bool {
	for j := i - 1; j >= 0; j-- {
		if tokens[j].typ == sqlWhitespace {
			continue
		}
		if tokens[j].typ == sqlWord {
			return strings.EqualFold(tokens[j].value, keyword)
		}
		return false
	}
	return false
}
