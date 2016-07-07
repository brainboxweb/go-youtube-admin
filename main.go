package main

import (
	"bytes"
	"code.google.com/p/google-api-go-client/youtube/v3"
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"text/template"
	"time"
)

func main() {

	app := cli.NewApp()
	app.Name = "youtube2"
	app.Usage = "Manage YouTube videos"
	app.Action = func(c *cli.Context) error {

		if c.NArg() > 0 {
			if c.Args().Get(0) == "backup" {
				backup()
			}
		}

		return nil
	}

	app.Run(os.Args)
}

func backup() {

	vids := getYouTubeData()

	d, err := yaml.Marshal(&vids)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	d1 := []byte(d)

	t := time.Now()
	filename := fmt.Sprintf("backup/youtube-%d-%d-%d.yml", t.Year(), t.Month(), t.Day())

	err = ioutil.WriteFile(filename, d1, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func getYouTubeData() []YouTubeData {

	service := getService()

	call := service.Channels.List("contentDetails").Mine(true)

	response, err := call.Do()
	if err != nil {
		// The channels.list method call returned an error.
		log.Fatalf("Error making API call to list channels: %v", err.Error())
	}

	data := []YouTubeData{}

	for _, channel := range response.Items {
		playlistId := channel.ContentDetails.RelatedPlaylists.Uploads
		//// Print the playlist ID for the list of uploaded videos.
		//fmt.Printf("Videos in list %s\r\n", playlistId)

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

				yt := YouTubeData{
					Title: playlistItem.Snippet.Title,
					Id:    playlistItem.Snippet.ResourceId.VideoId,
					Body:  playlistItem.Snippet.Description,
				}

				data = append(data, yt)

				//Get details!
				//Update the descriotiont
				//playlistItem.Snippet.Title += "\n+++"
			}

			// Set the token to retrieve the next page of results
			// or exit the loop if all results have been retrieved.
			nextPageToken = playlistResponse.NextPageToken
			if nextPageToken == "" {
				break
			}
		}
	}

	return data

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

func getPosts(postsFle string) map[int]Post {

	data := readYAMLFile(postsFle)
	posts := convertYAML(data)

	return posts
}

func readYAMLFile(filename string) []byte {

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatalf("Failed to read YML file : %v", err.Error())

	}

	return data
}

type YouTubeData struct {
	Id    string
	Title string
	Body  string
	Music []string
}

type Post struct {
	Title       string
	Description string
	Slug        string
	Date        string
	YouTubeData YouTubeData
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

const templateYouTube = `{{.YouTubeData.Body}}


_________________

"Development That Pays" is a weekly video that takes a business-focussed look at what's working now in software development. If you business depends on software development, we'd love to have you subscribe and join us!

SUBSCRIBE!
-- http://www.developmentthatpays.com/-/subscribe

LET'S CONNECT!
-- https://www.facebook.com/DevelopmentThatPays/
-- https://twitter.com/DevThatPays

{{if .YouTubeData.Music}}MUSIC{{ range .YouTubeData.Music }}
-- {{ . }}{{ end }}
{{ end }}
`
