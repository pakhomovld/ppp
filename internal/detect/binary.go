package detect

// isBinary checks whether the sample contains null bytes or a high ratio of
// non-text bytes, indicating binary content.
func isBinary(sample []byte) bool {
	// Check first 512 bytes (enough to detect binary reliably).
	n := len(sample)
	if n > 512 {
		n = 512
	}

	for i := 0; i < n; i++ {
		if sample[i] == 0 {
			return true
		}
	}
	return false
}
