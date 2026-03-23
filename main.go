package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/pakhomovld/pp/cmd"
	ppcolor "github.com/pakhomovld/pp/internal/color"
	"github.com/pakhomovld/pp/internal/detect"
	"github.com/pakhomovld/pp/internal/format"
	"github.com/pakhomovld/pp/internal/sniff"
)

func main() {
	cfg := cmd.ParseFlags()

	if cfg.Version {
		fmt.Printf("pp version %s\n", cmd.Version)
		os.Exit(0)
	}

	if cfg.NoColor {
		color.NoColor = true
	}

	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "pp: %v\n", err)
		os.Exit(1)
	}
}

func run(cfg cmd.Config) error {
	sr, err := sniff.NewReader(os.Stdin, sniff.DefaultSize)
	if err != nil {
		return err
	}

	sample := sr.Sample()
	if len(sample) == 0 {
		return nil
	}

	var detected detect.Format
	if cfg.ForceFormat != "" {
		detected = detect.Format(cfg.ForceFormat)
	} else {
		detected = detect.Detect(sample)
	}

	formatter := format.ForFormat(detected)

	var theme *ppcolor.Theme
	if ppcolor.ShouldColor() && !cfg.NoColor {
		theme = ppcolor.DefaultTheme()
	}

	return formatter.Format(os.Stdout, sr.Reader(), theme)
}
