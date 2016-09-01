package main

import (
	"code.google.com/p/google-api-go-client/youtube/v3"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestUnmarshalYAML(t *testing.T) {

	input := `
1-the-slug:
    title: The title
    description: "The description"
    date: 2015-08-20
    youtubedata:
        id: JkVr2DJM3Ac
        body: |-
            The body for YouTube purposes
        music:
            - "260809 Funky Nurykabe: http://ccmixter.org/files/jlbrock44/29186"
    body: |-
        This is the body
    transcript: |-
        This is the transcript.

2-the-slug-2:
    title: The title two
    description: "The description two"
    date: 2015-08-27
    youtubedata:
        id: xxxxxxxx
        body: |-
            The body for YouTube purposes. Again.
    body: |-
        This is the body. Again.
    transcript: |-
        This is the transcript. Again.`

	yt1 := YouTubeData{
		Id:    "JkVr2DJM3Ac",
		Body:  "The body for YouTube purposes",
		Music: []string{"260809 Funky Nurykabe: http://ccmixter.org/files/jlbrock44/29186"},
	}
	post1 := Post{
		Title:       "The title",
		Description: "The description",
		Date:        "2015-08-20",
		YouTubeData: yt1,
		Body:        "This is the body",
		Transcript:  "This is the transcript.",
	}

	yt2 := YouTubeData{
		Id:   "xxxxxxxx",
		Body: "The body for YouTube purposes. Again.",
	}
	post2 := Post{
		Title:       "The title two",
		Description: "The description two",
		Date:        "2015-08-27",
		YouTubeData: yt2,
		Body:        "This is the body. Again.",
		Transcript:  "This is the transcript. Again.",
	}

	expected := map[string]Post{}

	expected["1-the-slug"] = post1
	expected["2-the-slug-2"] = post2

	actual := convertYAML([]byte(input))

	eq := reflect.DeepEqual(expected["1-the-slug"], actual["1-the-slug"])
	if !eq {
		t.Errorf("expected %s, \n actual %s", expected["1-the-slug"], actual["1-the-slug"])
	}

	eq = reflect.DeepEqual(expected["2-the-slug-2"], actual["2-the-slug-2"])
	if !eq {
		t.Errorf("expected %s, \n actual %s", expected["2-the-slug-2"], actual["2-the-slug-2"])
	}
}

func TestReadYAMLFile(t *testing.T) {

	data := readYAMLFile("data/posts-test.yml")

	if data == nil {
		t.Error("Failed to read YAML file")
	}

}

func TestGetPosts(t *testing.T) {

	getPosts("data/posts.yml") //Just a test for parsing
}

func TestParseTemplate(t *testing.T) {

	post := Post{
		Title:       "The title",
		Description: "The description.",
		Date:        "2015-08-20",
		YouTubeData: YouTubeData{
			Id: "JkVr2DJM3Ac",
			Body: `The body for YouTube purposes.

On more than one line if necessary.`,
			Music: []string{
				"260809 Funky Nurykabe: http://ccmixter.org/files/jlbrock44/29186",
				"260809 Funky Nurykabe: http://ccmixter.org/files/jlbrock44/29186",
			},
		},
		Body:       "This is the body",
		Transcript: "This is the transcript.",
	}

	actual := parseTemplate(post)

	expected := parsed_1

	if actual != expected {
		t.Errorf("expected:\n %s, \n\n\n actual:\n %s", expected, actual)
	}

}

func TestGetVideo(t *testing.T) {

	video := getVideo("EHoyDH1cYwM")

	if !strings.Contains(video.Snippet.Title, "Jira") {
		t.Errorf("Video title does not contain 'Jira'")
	}
}

type FakeYouTube struct {
	Err error
}

func (yt FakeYouTube) persistVideo(*youtube.Video) error {

	return yt.Err
}

//
//func TestUpdateSnippet(t *testing.T) {
//
//	post := Post{
//		Title: "This is the Title of the Post",
//		YouTubeData: YouTubeData{
//			Id:    "EHoyDH1cYwM",
//			Title: "The original Youtube title",
//			Body:  "Thsi si the body om the post/youtube item",
//		},
//		Body: "this is the body of the POST item",
//	}
//
//	updateSnippet()
//}

func TestUpdateVideo(t *testing.T) {

	post := Post{
		Title: "This is the Title of the Post",
		YouTubeData: YouTubeData{
			Id:    "EHoyDH1cYwM",
			Title: "The original Youtube title",
			Body:  "Thsi si the body om the post/youtube item",
		},
		Body: "this is the body of the POST item",
	}

	c := make(chan interface{})

	yt := FakeYouTube{}

	go updateVideo(c, yt, 1, post)

	result := <-c

	//Assert the error
	err, found := result.(error)
	if found {
		t.Error("Video not updated", err.Error())
	}

}

func TestUpdateVideoErrorCondition(t *testing.T) {

	post := Post{
		Title: "This is the Title of the Post",
		YouTubeData: YouTubeData{
			Id:    "EHoyDH1cYwM",
			Title: "The original Youtube title",
		},
	}

	yt := FakeYouTube{
		Err: errors.New("Call to YoutTube Failed"),
	}

	c := make(chan interface{})

	go updateVideo(c, yt, 1, post)

	result := <-c

	//Assert the error
	_, found := result.(error)
	if !found {
		t.Errorf("YouTube error expected")
	}

}

const parsed_1 = `http://www.developmentthatpays.com The body for YouTube purposes.

On more than one line if necessary.


_________________

"Development That Pays" is a weekly video that takes a business-focused look at what's working now in Software Development.

If your business depends on Software Development, I'd love to have you subscribe for a new video every Wednesday!

SUBSCRIBE!
-- http://www.developmentthatpays.com/-/subscribe

LET'S CONNECT!
-- Facebook: https://www.facebook.com/DevelopmentThatPays/
-- Twitter: https://twitter.com/DevThatPays
-- LinkedIn: https://uk.linkedin.com/in/garystraughan

_________________



MUSIC
-- 260809 Funky Nurykabe: http://ccmixter.org/files/jlbrock44/29186
-- 260809 Funky Nurykabe: http://ccmixter.org/files/jlbrock44/29186

`
