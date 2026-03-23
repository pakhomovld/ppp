package detect

// Confidence represents how certain a detector is about the format.
type Confidence int

const (
	None Confidence = iota
	Low
	Medium
	High
)

// String returns the confidence level as a lowercase string.
func (c Confidence) String() string {
	switch c {
	case None:
		return "none"
	case Low:
		return "low"
	case Medium:
		return "medium"
	case High:
		return "high"
	default:
		return "none"
	}
}
