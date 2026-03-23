package cmd

import (
	"flag"
	"fmt"
	"os"
)

// Version is set via ldflags at build time.
var Version = "dev"

// Config holds CLI flags.
type Config struct {
	ForceFormat string
	NoColor     bool
	Version     bool
	Inspect     bool
	Strict      bool
}

const helpText = `ppp — universal pipe pretty-printer

Usage:
  <command> | ppp [flags]

Reads from stdin, auto-detects the data format, and pretty-prints
it with syntax highlighting. No flags needed.

Flags:
  -f, --format <fmt>   Force a specific format (skip auto-detection)
      --inspect        Output detection metadata as JSON, then exit
      --strict         Exit 2 if detection confidence is low or none
      --no-color       Disable colored output
  -v, --version        Print version
  -h, --help           Print this help

Supported formats:
  json        JSON objects and arrays
  yaml        YAML documents
  csv         Comma-separated values
  tsv         Tab-separated values
  xml         XML documents
  html        HTML documents
  toml        TOML configuration
  jwt         JSON Web Tokens (decoded header + payload)
  base64      Base64-encoded data (decoded, inner format detected)
  urlencoded  URL-encoded key=value pairs
  log         Log lines with timestamps and levels

Examples:
  curl -s https://api.example.com/data | ppp
  cat config.yaml | ppp
  echo '{"name":"Alice"}' | ppp
  docker logs myapp | ppp
  echo 'eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxIn0.sig' | ppp
  cat data.csv | ppp -f csv

Exit codes:
  0   Success
  1   I/O or format error
  2   Low confidence (only with --strict)

Color is automatically disabled when stdout is not a terminal
or the NO_COLOR environment variable is set.

More info: https://github.com/pakhomovld/ppp`

// ParseFlags parses CLI arguments.
func ParseFlags() Config {
	var cfg Config

	flag.StringVar(&cfg.ForceFormat, "format", "", "force a specific format (json, yaml, csv, ...)")
	flag.StringVar(&cfg.ForceFormat, "f", "", "force a specific format (shorthand)")
	flag.BoolVar(&cfg.NoColor, "no-color", false, "disable colored output")
	flag.BoolVar(&cfg.Inspect, "inspect", false, "output detection metadata as JSON")
	flag.BoolVar(&cfg.Strict, "strict", false, "exit 2 if confidence is low or none")
	flag.BoolVar(&cfg.Version, "version", false, "print version")
	flag.BoolVar(&cfg.Version, "v", false, "print version (shorthand)")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, helpText)
	}

	flag.Parse()
	return cfg
}
