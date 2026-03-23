# ppp

Universal pipe pretty-printer. One command, any format.

```
<command> | ppp
```

`ppp` reads from stdin, auto-detects the data format, and pretty-prints it with syntax highlighting. No flags needed.

![ppp demo](demo/hero.gif)

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
brew install pakhomovld/tap/ppp
```

**Go:**

```sh
go install github.com/pakhomovld/pp@latest
```

The module path is `pp`; the binary is `ppp` (one more `p` for *pretty*).

**Binary:** Download from [GitHub Releases](https://github.com/pakhomovld/ppp/releases).

## Examples

**JSON** — auto-detected, indented, colored:

![JSON demo](demo/json.gif)

**CSV** — auto-detected, rendered as aligned table:

![CSV demo](demo/csv.gif)

**JWT** — decoded header and payload:

![JWT demo](demo/jwt.gif)

**Logs** — timestamps dimmed, levels colorized:

![Logs demo](demo/logs.gif)

**Base64** — decoded, then inner format detected:

```sh
echo '{"secret":"found"}' | base64 | ppp
```

```
[base64 decoded → json]
{
  "secret": "found"
}
```

**URL-encoded** — decoded and displayed as table:

```sh
echo 'user=john%40example.com&token=abc&active=true' | ppp
```

```
  active = true
  token  = abc
  user   = john@example.com
```

## Real-World Usage

```sh
kubectl get pods -o yaml | ppp
curl -s https://api.github.com/users/octocat | ppp
journalctl -u nginx -n 50 | ppp
echo $JWT_TOKEN | ppp
pbpaste | ppp
cat mystery-file.txt | ppp --inspect
curl -s api.example.com/data | ppp --strict
```

## Flags

```
--format, -f <fmt>   Force a specific format (skip auto-detection)
--inspect            Output detection metadata as JSON, then exit
--strict             Exit 2 if detection confidence is low or none
--no-color           Disable colored output
--version, -v        Print version
--help, -h           Print help
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | I/O or format error |
| 2 | Low confidence (only with `--strict`) |

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

`ppp` replaces all of them with one command. You don't need to know the format — just pipe to `ppp`.

## How Detection Works

`ppp` reads the first 8KB of stdin and runs format detectors in priority order:

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

## Detection Notes

| Input | Detected As | Why |
|-------|-------------|-----|
| `"Dear John: hope you are well"` | Plain | Single `key: value` is too ambiguous for YAML |
| `"<root><unclosed"` | XML (medium) | Starts with `<`, partial tags detected |
| Long English words | Plain | Not enough base64 signals (min 20 chars + decodable) |
| Binary / null bytes | Plain | Non-UTF-8 content falls through |
| Whitespace-only | Plain | Empty after trim |

Use `--inspect` to see what `ppp` detected and how confident it is.

## Build from Source

```sh
git clone https://github.com/pakhomovld/ppp
cd ppp
make build     # produces ./ppp
make test      # runs all tests
make install   # copies to /usr/local/bin
```

## License

MIT
