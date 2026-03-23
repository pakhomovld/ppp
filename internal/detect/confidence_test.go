package detect

import "testing"

func TestConfidence_String(t *testing.T) {
	tests := []struct {
		c    Confidence
		want string
	}{
		{None, "none"},
		{Low, "low"},
		{Medium, "medium"},
		{High, "high"},
		{Confidence(99), "none"},
	}

	for _, tt := range tests {
		got := tt.c.String()
		if got != tt.want {
			t.Errorf("Confidence(%d).String() = %q, want %q", tt.c, got, tt.want)
		}
	}
}
