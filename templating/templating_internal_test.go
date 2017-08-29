package templating

import (
	"fmt"
	"strings"
	"testing"
)

func TestTemplateNotFound(t *testing.T) {
	data := YouTubeData{
		Id:         "sdsdsdsd",
		Title:      "The TitleTitle",
		Body:       "the Body",
		Transcript: "the transcript",
		TopResult:  "http://www.top.com",
		Music:      []string{},
	}
	_, err := applyTemplate(data, "404.txt")
	if err == nil {
		t.Errorf("expected err. Not founc")
	}
}

func TestParse(t *testing.T) {
	id := "ididididid"
	title := "The Video Title"
	topResult := "http://number-one-on-google.com"
	templateFile := "youtube.txt"

	data := YouTubeData{
		Id:              id,
		Title:           title,
		BodyFirst:       "The Body",
		BodyAllButFirst: "Second line of the body",
		TopResult:       topResult,
		Music:           []string{},
	}

	parsed, _ := parseTemplate(data, templateFile)

	testCases := []struct {
		expected string
	}{
		{"The Video Title // "},
		{" // The Body"},
		//{"http://www.DevelopmentThatPays.com/-/subscribe"},
		{"https://www.youtube.com/watch?v=" + id},
		{topResult},
	}

	for _, tc := range testCases {
		t.Run("----", func(t *testing.T) {
			if !strings.Contains(parsed, tc.expected) {
				t.Errorf("String '%s' not found.", tc.expected)
			}
		})
	}
}

func TestSplitBody(t *testing.T) {
	testCases := []struct {
		input     string
		expected1 string
		expected2 string
	}{
		{
			`Now is the time

For all good men`,
			"Now is the time",
			"For all good men",
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.input), func(t *testing.T) {
			actual1, actual2 := splitBody(tc.input)

			if actual1 != tc.expected1 {
				t.Errorf("got %s; want %s", actual1, tc.expected1)
			}
			if actual2 != tc.expected2 {
				t.Errorf("got %s; want %s", actual2, tc.expected2)
			}
		})
	}
}

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
