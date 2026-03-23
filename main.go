package main

import (
	"encoding/json"
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
		fmt.Printf("ppp version %s\n", cmd.Version)
		os.Exit(0)
	}

	if cfg.NoColor {
		color.NoColor = true
	}

	exitCode, err := run(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ppp: %v\n", err)
		os.Exit(1)
	}
	os.Exit(exitCode)
}

func run(cfg cmd.Config) (int, error) {
	sr, err := sniff.NewReader(os.Stdin, sniff.DefaultSize)
	if err != nil {
		return 1, err
	}

	sample := sr.Sample()
	if len(sample) == 0 {
		return 0, nil
	}

	var result detect.Result
	if cfg.ForceFormat != "" {
		result = detect.Result{Format: detect.Format(cfg.ForceFormat), Confidence: detect.High}
	} else {
		result = detect.Detect(sample)
	}

	lowConfidence := result.Confidence <= detect.Low

	if cfg.Inspect {
		meta := struct {
			Format     string `json:"format"`
			Confidence string `json:"confidence"`
		}{
			Format:     string(result.Format),
			Confidence: result.Confidence.String(),
		}
		enc := json.NewEncoder(os.Stdout)
		if err := enc.Encode(meta); err != nil {
			return 1, err
		}
		if cfg.Strict && lowConfidence {
			return 2, nil
		}
		return 0, nil
	}

	if cfg.Strict && lowConfidence {
		fmt.Fprintf(os.Stderr, "ppp: low confidence detection (%s: %s)\n", result.Format, result.Confidence)
		return 2, nil
	}

	formatter := format.ForFormat(result.Format)

	var theme *ppcolor.Theme
	if ppcolor.ShouldColor() && !cfg.NoColor {
		theme = ppcolor.DefaultTheme()
	}

	if err := formatter.Format(os.Stdout, sr.Reader(), theme); err != nil {
		return 1, err
	}
	return 0, nil
}
