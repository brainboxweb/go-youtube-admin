package templating

import (
	"fmt"
	"testing"
)

func TestTruncate(t *testing.T) {
	expectedLength := 4
	s := "123456789"
	output := truncate(s, expectedLength)
	actualLength := len(output)
	if actualLength != expectedLength {
		t.Errorf("expected %d. Got %d", expectedLength, actualLength)
	}
}

func TestRemoveLinebreaks(t *testing.T) {
	s := `one
two



three`
	expected := "one two three"
	actual := removeLineBreaks(s)
	if actual != expected {
		t.Errorf("expected %s. Got %s", expected, actual)
	}
}

func TestStripMarkdown(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"This is **bold**", "This is bold"},
		{"[What fun!](http://thinng.com)", "What fun!"},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.input), func(t *testing.T) {
			actual := stripMarkdown(tc.input)

			if actual != tc.expected {
				t.Errorf("got %s; want %s", actual, tc.expected)
			}
		})
	}
}
