package color

import "testing"

func TestGet(t *testing.T) {
	tests := []struct {
		value    string
		expected string
	}{
		{value: "red", expected: Red},
		{value: "green", expected: Green},
		{value: "yellow", expected: Yellow},
		{value: "blue", expected: Blue},
		{value: "magenta", expected: Magenta},
		{value: "cyan", expected: Cyan},
		{value: "gray", expected: Gray},
		{value: "white", expected: White},
		{value: "reset", expected: Reset},
		{value: "RED", expected: Red},
		{value: "Green", expected: Green},
		{value: "unknown", expected: Reset},
		{value: "", expected: Reset},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if result := Get(tt.value); result != tt.expected {
				t.Errorf("Get(%q) = %q, expected %q", tt.value, result, tt.expected)
			}
		})
	}
}
