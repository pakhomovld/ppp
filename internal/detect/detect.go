package detect

// Format identifies a data format.
type Format string

const (
	Binary    Format = "binary"
	NDJSON    Format = "ndjson"
	JSON      Format = "json"
	YAML      Format = "yaml"
	CSV       Format = "csv"
	TSV       Format = "tsv"
	TOML      Format = "toml"
	XML       Format = "xml"
	HTML      Format = "html"
	JWT       Format = "jwt"
	Base64    Format = "base64"
	URLEncode Format = "urlencoded"
	LogLine   Format = "log"
	Markdown  Format = "markdown"
	Plain     Format = "plain"
)

// Result is returned by each detector.
type Result struct {
	Format     Format
	Confidence Confidence
}

// Detector inspects a sample and returns a detection result.
type Detector interface {
	Detect(sample []byte) Result
}

// Detect runs all detectors in priority order and returns the best match.
func Detect(sample []byte) Result {
	if len(sample) == 0 {
		return Result{Format: Plain, Confidence: None}
	}

	if isBinary(sample) {
		return Result{Format: Binary, Confidence: High}
	}

	detectors := []Detector{
		&JWTDetector{},
		&NDJSONDetector{},
		&JSONDetector{},
		&XMLDetector{},
		&YAMLDetector{},
		&TOMLDetector{},
		&CSVDetector{},
		&URLDetector{},
		&LogDetector{},
		&Base64Detector{},
	}

	best := Result{Format: Plain, Confidence: None}
	for _, d := range detectors {
		r := d.Detect(sample)
		if r.Confidence > best.Confidence {
			best = r
		}
		if best.Confidence == High {
			break
		}
	}
	return best
}
