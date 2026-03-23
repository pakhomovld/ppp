# pp

Universal pipe pretty-printer. One command, any format.

```
<command> | pp
```

`pp` reads from stdin, auto-detects the data format, and pretty-prints it with syntax highlighting. No flags needed.

## Supported Formats

| Format | Detection | Output |
|--------|-----------|--------|
| JSON | `{` or `[` prefix + valid parse | Indented + colored |
| YAML | `---` separator or `key: value` lines | Colored keys/values |
| CSV/TSV | Consistent delimiter counts | Aligned table |
| XML | `<` prefix, `<?xml` declaration | Indented + colored tags |
| HTML | `<!DOCTYPE html>` or `<html>` | Indented + colored |
| TOML | `[section]` + `key = value` | Colored sections/keys |
| Log lines | Timestamps + level patterns | Colorized levels |
| JWT | Three base64url segments with `"alg"` header | Decoded header + payload |
| Base64 | Valid charset + decodable to UTF-8 | Decoded + inner format detected |
| URL-encoded | `key=value&key=value` | Decoded key-value table |

Unknown formats pass through unchanged.

## Install

**Homebrew:**

```sh
brew install pakhomovld/tap/pp
```

**Go:**

```sh
go install github.com/pakhomovld/pp@latest
```

**Binary:** Download from [GitHub Releases](https://github.com/pakhomovld/pp/releases).

## Examples

**JSON** — auto-detected, indented, colored:

```sh
curl -s https://api.example.com/users | pp
```

```
{
  "name": "Alice",
  "active": true,
  "age": 30
}
```

**CSV** — auto-detected, rendered as aligned table:

```sh
cat data.csv | pp
```

```
  name     age  city
  ──────────────────
  Alice    30   NYC
  Bob      25   LA
```

**JWT** — decoded header and payload:

```sh
echo 'eyJhbGciOiJIUzI1NiJ9.eyJuYW1lIjoiSm9obiJ9.signature' | pp
```

```
Header:
{
  "alg": "HS256"
}

Payload:
{
  "name": "John"
}

Signature: [not verified]
```

**Logs** — timestamps dimmed, levels colorized:

```sh
docker logs myapp | pp
```

**Base64** — decoded, then inner format detected:

```sh
echo '{"secret":"found"}' | base64 | pp
```

```
[base64 decoded → json]
{
  "secret": "found"
}
```

**URL-encoded** — decoded and displayed as table:

```sh
echo 'user=john%40example.com&token=abc&active=true' | pp
```

```
  active = true
  token  = abc
  user   = john@example.com
```

## Flags

```
--format, -f <fmt>   Force a specific format (skip auto-detection)
--no-color           Disable colored output
--version, -v        Print version
--help, -h           Print help
```

Color is automatically disabled when:
- stdout is not a terminal (piping to a file or another command)
- `NO_COLOR` environment variable is set

## Why not jq / yq / bat?

| Tool | Limitation |
|------|-----------|
| `jq` | JSON only. Must know it's JSON before piping. |
| `yq` | YAML only. |
| `bat` | Syntax highlighting only — doesn't restructure, reindent, or decode. |
| `column -t` | Tables only. No auto-detection. |
| `xmllint` | XML only. No color. |
| `base64 -d` | Decodes but doesn't format the result. |

`pp` replaces all of them with one command. You don't need to know the format — just pipe to `pp`.

## How Detection Works

`pp` reads the first 8KB of stdin and runs format detectors in priority order:

1. **JWT** — three dot-separated base64url segments
2. **JSON** — starts with `{` or `[`, valid parse
3. **XML/HTML** — starts with `<`
4. **YAML** — `---` or `key: value` patterns
5. **TOML** — `[section]` headers + `key = value`
6. **CSV/TSV** — consistent delimiter counts across lines
7. **URL-encoded** — `key=value&key=value` pattern
8. **Log lines** — timestamp and level patterns
9. **Base64** — valid charset, decodes to UTF-8
10. **Plain text** — fallback, passes through unchanged

Each detector returns a confidence score (High/Medium/Low/None). Highest confidence wins. On ties, earlier detector wins.

## Build from Source

```sh
git clone https://github.com/pakhomovld/pp
cd pp
make build     # produces ./pp
make test      # runs all tests
make install   # copies to /usr/local/bin
```

## License

MIT
