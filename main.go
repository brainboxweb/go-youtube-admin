package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"code.google.com/p/google-api-go-client/youtube/v3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	"github.com/brainboxweb/go-youtube-admin/templating"
)

const database = "../go-posts-admin/db/dtp.db"
const templateFile = "templating/youtube.txt"

func main() {
	app := cli.NewApp()
	app.Name = "youtube2"
	app.Usage = "Manage YouTube videos"
	app.Action = func(c *cli.Context) error {
		if c.NArg() > 0 {
			if c.Args().Get(0) == "backup" {
				backup()
			}
			if c.Args().Get(0) == "update" {

				update(c.Args().Get(1))
			}
		}
		return nil
	}
	app.Run(os.Args)
}

//YouTuber is an interface for the YouTube client
type YouTuber interface {
	persistVideo(*youtube.Video) error
}

//MyYouTube is a struct representing the YT service
type MyYouTube struct{} //@todo - rename this
//Implement the YouTuber interface
func (MyYouTube) persistVideo(video *youtube.Video) error {
	service := getService()
	call := service.Videos.Update("snippet", video)
	_, err := call.Do()
	return err
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

func update(id string) {
	yt := MyYouTube{}
	posts := getPosts(id)
	c := make(chan UpdateResult)
	for id, post := range posts {
		go updateVideo(c, yt, id, post)
	}
	for i := 0; i < len(posts); i++ {
		result := <-c
		fmt.Println(result)
	}
}

//UpdateResult provides information on a video update
type UpdateResult struct{
	Status string
	Error error
}

func updateVideo(c chan UpdateResult, yt YouTuber, index int, post Post) {
	videoID := post.YouTubeData.Id
	video := getVideo(videoID)
	updated := updateSnippet(video, index, post)

	if !updated {
		c <- UpdateResult{
			Status: fmt.Sprintf("NO CHANGE - %d %s", index, post.Title),
			Error: nil,
		}
		return
	}
	err := yt.persistVideo(video)
	if err != nil {
		c <- UpdateResult{
			Status: "",	
			Error: fmt.Errorf(">>>>ERROR - %d %s, %s", index, post.Title, err),
		}
	}
	c <- UpdateResult{
		Status: fmt.Sprintf("UPDATED - %d %s", index, post.Title),
		Error: nil,
	} 
	
}

func getVideo(videoID string) *youtube.Video {
	service := getService()
	call := service.Videos.List("snippet").Id(videoID)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making API call to get video  data: %v", err.Error()) // The channels.list method call returned an error.
	}
	if len(response.Items) < 1 {
		panic("Video id not found:" + videoID)
	}
	return response.Items[0]
}

func updateSnippet(video *youtube.Video, index int, post Post) (updated bool) {

	updated = false

	//Title
	newTitle := post.Title
	if video.Snippet.Title != newTitle {
		updated = true
		video.Snippet.Title = newTitle
	}

	//Tags
	commonTags := []string{
		"Development That Pays",
	}
	post.Keywords = append(post.Keywords, commonTags...)
	if compareSlice(video.Snippet.Tags, post.Keywords) == false {
		updated = true
		video.Snippet.Tags = post.Keywords
	}

	//Description
	data := templating.YouTubeData{
		Id:           post.YouTubeData.Id,
		Playlist:     post.YouTubeData.Playlist,
		Index:        index,
		Title:        post.Title,
		Description:  post.Description,
		Body:         post.YouTubeData.Body,
		Transcript:   post.Transcript,
		TopResult:    post.TopResult,
		Music:        post.YouTubeData.Music,
		ClickToTweet: post.ClickToTweet,
	}
	newDescription, err := templating.GetYouTubeBody(data, templateFile)
	if err != nil {
		panic(err)
	}
	if video.Snippet.Description != newDescription {
		video.Snippet.Description = newDescription
		updated = true
	}
	return updated
}

func compareSlice(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
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
					Id:   playlistItem.Snippet.ResourceId.VideoId,
					Body: playlistItem.Snippet.Description,
				}
				data = append(data, yt)
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

func getService() (service *youtube.Service) {
	//@TODO - ADD SYNC.ONCE
	var once sync.Once
	once.Do(func() {
		client, err := buildOAuthHTTPClient(youtube.YoutubeScope)
		if err != nil {
			log.Fatalf("Error building OAuth client: %v", err)
		}
		service, err = youtube.New(client)
		if err != nil {
			log.Fatalf("Error creating YouTube client: %v", err)
		}
	})
	return service
}

func getPosts(id string) map[int]Post {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	posts := make(map[int]Post)
	q := "SELECT posts.id, slug, title, description, topresult, click_to_tweet, transcript, youtube.id AS youtube_id, youtube.body AS youtube_body, coalesce(youtube.playlist, '') AS youtube_playlist FROM posts LEFT JOIN youtube ON posts.id = youtube.post_id"
	if id != "" {
		q = q + fmt.Sprintf(" WHERE posts.id = %s", id)
	}

	rows, err := db.Query(q)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		p := new(Post)
		err = rows.Scan(&p.Id, &p.Slug, &p.Title, &p.Description, &p.TopResult, &p.ClickToTweet, &p.Transcript, &p.YouTubeData.Id, &p.YouTubeData.Body, &p.YouTubeData.Playlist)
		if err != nil {
			panic(err)
		}
		rows2, err := db.Query("SELECT keyword_id FROM posts_keywords_xref WHERE post_id = ? ORDER BY sort_order", p.Id)
		if err != nil {
			panic(err)
		}
		for rows2.Next() {
			keyword := ""
			err = rows2.Scan(&keyword)
			if err != nil {
				panic(err)
			}
			p.Keywords = append(p.Keywords, keyword)
		}
		//Music
		rows4, err := db.Query("SELECT music_id FROM youtube_music_xref WHERE youtube_id = ?", p.YouTubeData.Id)
		if err != nil {
			panic(err)
		}
		for rows4.Next() {
			var music string
			err = rows4.Scan(&music)
			if err != nil {
				panic(err)
			}
			p.YouTubeData.Music = append(p.YouTubeData.Music, music)
		}
		posts[p.Id] = *p
	}
	return posts
}

// YouTubeData is a a container for video data
type YouTubeData struct {
	Id       string
	Body     string
	Playlist string
	Music    []string
}

//Post is a container for submitted data relating to a video
type Post struct {
	Id           int
	Slug         string
	Title        string
	Description  string
	Date         string
	TopResult    string
	Keywords     []string
	YouTubeData  YouTubeData
	Body         string
	Transcript   string
	ClickToTweet string
}
