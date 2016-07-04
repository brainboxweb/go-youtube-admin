package main

import (
	"flag"
	//"fmt"
	"log"

	"code.google.com/p/google-api-go-client/youtube/v3"

	//"net/http"
	"fmt"
	//"regexp"
	//"gopkg.in/yaml.v2"

	//"io/ioutil"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	//"sync"
	"text/template"
	//"bytes"
	"bytes"
)

func main() {
	flag.Parse()

	service := getService()

	call := service.Channels.List("contentDetails").Mine(true)

	response, err := call.Do()
	if err != nil {
		// The channels.list method call returned an error.
		log.Fatalf("Error making API call to list channels: %v", err.Error())
	}

	for _, channel := range response.Items {
		playlistId := channel.ContentDetails.RelatedPlaylists.Uploads
		// Print the playlist ID for the list of uploaded videos.
		fmt.Printf("Videos in list %s\r\n", playlistId)

		nextPageToken := ""
		for {
			// Call the playlistItems.list method to retrieve the
			// list of uploaded videos. Each request retrieves 50
			// videos until all videos have been retrieved.
			playlistCall := service.PlaylistItems.List("snippet").
				PlaylistId(playlistId).
				MaxResults(50).
				PageToken(nextPageToken)

			playlistResponse, err := playlistCall.Do()

			if err != nil {
				// The playlistItems.list method call returned an error.
				log.Fatalf("Error fetching playlist items: %v", err.Error())
			}

			for _, playlistItem := range playlistResponse.Items {

				title := playlistItem.Snippet.Title
				videoId := playlistItem.Snippet.ResourceId.VideoId

				fmt.Printf("%v, (%v)\r\n", title, videoId)

				//Get details!
				//Update the descriotiont

				playlistItem.Snippet.Title += "\n+++"

			}

			// Set the token to retrieve the next page of results
			// or exit the loop if all results have been retrieved.
			nextPageToken = playlistResponse.NextPageToken
			if nextPageToken == "" {
				break
			}
		}
	}

}

func getService() *youtube.Service {

	client, err := buildOAuthHTTPClient(youtube.YoutubeScope)
	if err != nil {
		log.Fatalf("Error building OAuth client: %v", err)
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	return service
}

func updateSignoff(service *youtube.Service, id string) {

	//GET the video details
	call := service.Videos.List("snippet").Id(id)
	videoResponse, err := call.Do()
	if err != nil {
		// The channels.list method call returned an error.
		log.Fatalf("Error making API call to get video: %v", err.Error())
	}
	for _, videoItem := range videoResponse.Items {

		//description := videoItem.Snippet.Description

		videoItem.Snippet.Description += "+++"

		//Try an update
		call := service.Videos.Update("snippet", videoItem)
		_, err := call.Do()
		if err != nil {
			// The channels.list method call returned an error.
			log.Fatalf("Error making API call to UPDATE video: %v", err.Error())
		}

		return
	}

}

//
//func updateDescriptionFooter(description string) string {
//
//	r := regexp.MustCompile(`(?s:(.*)----(.*)?)`)
//	fmt.Println(r)
//
//	matches := r.FindStringSubmatch(description)
//
//	if len(matches) < 2 {
//		return description
//	}
//
//	//for _, match := range matches{
//	//
//	//	fmt.Println("\n\n=========================\n", match)
//	//
//	//}
//
//	out := matches[1] + "----\n" + "NEW CONTENT"
//
//	//out := r.ReplaceAllString(description, "FOOTER")
//
//	return out
//}

func getPosts(postsFle string) map[int]Post {

	data := readYAMLFile(postsFle)
	posts := convertYAML(data)
	fmt.Println(posts)

	return posts

}

func readYAMLFile(filename string) []byte {

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatalf("Failed to read YML file : %v", err.Error())

	}

	return data
}

type Post struct {
	Title       string
	Description string
	Slug        string
	Date        string `yaml:"date"`
	YouTubeId   string `yaml:youtubeid"`
	YouTubeBody string `yaml:youtubebody"`
	Image       string
	Body        string
	Transcript  string
}

func convertYAML(input []byte) map[int]Post {
	posts := make(map[int]Post)

	err := yaml.Unmarshal(input, &posts)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return posts
}

//type templateHandler struct {
//	once     sync.Once
//	filename string
//	templ    *template.Template
//}

func parseTemplate(post Post) string {

	t := template.New("Post template") // Create a template.
	t, err := t.Parse(templateYouTube) // Parse template file.
	if err != nil {
		log.Fatal("error: %v", err)
	}

	var buff bytes.Buffer

	t.Execute(&buff, post)

	return buff.String()
}

const templateYouTube = `{{.YouTubeBody}}


_________________

"Development That Pays" is a weekly video that takes a business-focussed look at what's working now in software development. If you business depends on software development, we'd love to have you subscribe and join us!

SUBSCRIBE!
-- http://www.developmentthatpays.com/-/subscribe

LET'S CONNECT!
-- https://www.facebook.com/DevelopmentThatPays/
-- https://twitter.com/DevThatPays

MUSIC
-- 260809 Funky Nurykabe: http://ccmixter.org/files/jlbrock44/29186`
