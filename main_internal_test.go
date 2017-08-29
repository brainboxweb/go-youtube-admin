package main

import (
	"code.google.com/p/google-api-go-client/youtube/v3"
	"strings"
	"testing"
)

func TestGetPosts(t *testing.T) {
	getPosts() //Just a test for parsing
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

func TestUpdateVideo(t *testing.T) {

	post := Post{
		Title: "This is the Title of the Post",
		YouTubeData: YouTubeData{
			Id:   "EHoyDH1cYwM",
			Body: "Thsi si the body om the post/youtube item",
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

//func TestUpdateVideoErrorCondition(t *testing.T) {
//
//	post := Post{
//		Title: "This is the Title of the Post",
//		YouTubeData: YouTubeData{
//			Id: "EHoyDH1cYwM",
//		},
//		Transcript: "Now is the time for all good men to come to the aid",
//	}
//
//	yt := FakeYouTube{
//		Err: errors.New("Call to YoutTube Failed"),
//	}
//
//	c := make(chan interface{})
//
//	go updateVideo(c, yt, 1, post)
//
//	result := <-c
//
//	//Assert the error
//	_, found := result.(error)
//
//	if !found {
//		t.Errorf("YouTube error expected")
//	}
//}
