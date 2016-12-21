package main

import (
	"code.google.com/p/google-api-go-client/youtube/v3"
	"fmt"
	//"github.com/brainboxweb/go-youtube-admin/bitly"
	"github.com/brainboxweb/go-youtube-admin/templating"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	//"net/http"
	//"net/url"
	"os"
	//"strconv"
	//"strings"
	"sync"
	"time"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const database = "../go-posts-admin/db/dtp.db"
//const postsFile = "../go-posts-utils/data/posts.yml"
//const tweetsFile = "data/tweets.yml"
const tagsFile = "data/tags.yml"

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
				update()
			}

			//if c.Args().Get(0) == "bitly" {
			//	updateBitly()
			//}

			//if c.Args().Get(0) == "getTags" {
			//	getTags()
			//}
		}

		return nil
	}

	app.Run(os.Args)
}

//YouTuber is an interface for the YouTube client
type YouTuber interface {
	persistVideo(*youtube.Video) error
}

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

func update() {

	//fetchRemotePostsData()

	yt := MyYouTube{}

	posts := getPosts()

	//tweets := getTweets(tweetsFile)

	c := make(chan interface{})

	for id, post := range posts {

		//parts := strings.Split(slug, "-")
		//k, err := strconv.Atoi(parts[0])
		//if err != nil {
		//	log.Fatalf("error: %v", err)
		//}

		//slug := fmt.Sprintf("%d-%s", id, post.Slug)

		go updateVideo(c, yt, id, post)
	}

	for i := 0; i < len(posts); i++ {
		result := <-c
		fmt.Println(result)
	}

}
//
//// This should move to ADMIN!
//func updateBitly() {
//
//	posts := getPosts(
//		postsFile)
//
//	tweets := getTweets(tweetsFile)
//
//	for slug, post := range posts {
//		if _, ok := tweets[slug]; ok {
//			continue
//		}
//		status := fmt.Sprintf("%s --> https://youtu.be/%s via @DevThatPays", post.Title, post.YouTubeData.Id)
//		var Url *url.URL
//		Url, err := url.Parse("http://twitter.com/home/")
//		if err != nil {
//			panic("boom")
//		}
//		parameters := url.Values{}
//		parameters.Add("status", status)
//		Url.RawQuery = parameters.Encode()
//
//		link := bitly.GetShortnedLink(Url.String())
//
//		tweets[slug] = Tweet{link}
//
//	}
//
//	data, err := yaml.Marshal(tweets)
//	if err != nil {
//		panic("did not see that coming")
//	}
//
//	err = ioutil.WriteFile(tweetsFile, data, 0644)
//	if err != nil {
//		panic("more surprises")
//	}
//
//}

type Tags struct {
	Tags []string
}

//func getTags() {
//
//	posts := getPosts()
//
//	tagsMap := make(map[int]Tags)
//
//	for slug, post := range posts {
//
//		video := getVideo(post.YouTubeData.Id)
//
//		tagsData := video.Snippet.Tags
//
//		tags := Tags{}
//		for _, tagData := range tagsData {
//			//tag :=  Tag{tagData}
//
//			tags.Tags = append(tags.Tags, tagData)
//		}
//
//		tagsMap[slug] = tags
//
//		fmt.Println("\n", slug)
//
//	}
//
//	data, err := yaml.Marshal(tagsMap)
//	if err != nil {
//		panic("did not see that coming")
//	}
//
//	err = ioutil.WriteFile(tagsFile, data, 0644)
//	if err != nil {
//		panic("more surprises")
//	}
//
//}
//
//func fetchRemotePostsData() {
//
//	client := &http.Client{}
//
//	resp, err := client.Get("http://www.developmentthatpays.com/" + postsFile)
//
//	if err != nil {
//		log.Fatalf("error: %v", err)
//	}
//
//	defer resp.Body.Close()
//
//	data, err := ioutil.ReadAll(resp.Body)
//
//	if err != nil {
//		log.Fatalf("error: %v", err)
//	}
//
//	err = ioutil.WriteFile(postsFile, data, 0644)
//	if err != nil {
//		log.Fatalf("error: %v", err)
//	}
//}

func updateVideo(c chan interface{}, yt YouTuber, index int, post Post) {

	videoId := post.YouTubeData.Id
	video := getVideo(videoId)

	updated := updateSnippet(video, index, post)
	if !updated {

		//video.Snippet.Tags
		c <- fmt.Sprintf("NO CHANGE - %d %s", index, post.Title)
		return
	}

	//err := yt.persistVideo(video)
	//if err != nil {
	//	//c <- fmt.Sprintf("ERROR - %d %s - %s", index, post.Title, err)
	//	c <- err //Not great - don't know the source of the error
	//}

	c <- fmt.Sprintf("UPDATED - %d %s %v", index, post.Title)
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

	newTitle := fmt.Sprintf("%s  - DTP#%d", post.Title, index)
	if video.Snippet.Title != newTitle {
		video.Snippet.Title = newTitle
		updated = true
	}

	data := templating.YouTubeData{
		Id:           post.YouTubeData.Id,
		Title:        post.Title,
		Body:         post.YouTubeData.Body,
		Transcript:   post.Transcript,
		TopResult:    post.TopResult,
		Music:        post.YouTubeData.Music,
		ClickToTweet: post.ClickToTweet,
	}

	newDescription := templating.GetYouTubeBody(data, templateFile)
	//if err != nil {
	//	panic("Error experienced when creating newDescription")
	//}

	if video.Snippet.Description != newDescription {
		video.Snippet.Description = newDescription
		updated = true
	}
	return updated
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
					Id:   playlistItem.Snippet.ResourceId.VideoId,
					Body: playlistItem.Snippet.Description,
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

func getPosts() map[int]Post {

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	posts := make(map[int]Post)

	rows, err := db.Query("SELECT id, slug, title, description, published, body, transcript, topresult, click_to_tweet FROM posts")
	if err != nil {
		panic(err)
	}

	for rows.Next() {


		p := new(Post)

		err = rows.Scan(&p.Id, &p.Slug, &p.Title, &p.Description, &p.Date, &p.Body, &p.Transcript, &p.TopResult, &p.ClickToTweet)
		if err != nil {
			panic(err)
		}

		//Tags/Keywords
		//keywords := new(Keywords)

		rows2, err := db.Query("SELECT keyword_id FROM posts_keywords_xref WHERE post_id = ?", p.Id)
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

		//Youtube
		yt := new(YouTubeData)

		rows3, err := db.Query("SELECT id, body FROM youtube WHERE post_id = ?", p.Id)
		if err != nil {
			panic(err)
		}

		for rows3.Next() {
			err = rows3.Scan(&yt.Id, &yt.Body)
			if err != nil {
				panic(err)
			}

			rows4, err := db.Query("SELECT music_id  FROM youtube_music_xref WHERE youtube_id = ?", yt.Id)
			if err != nil {
				panic(err)
			}

			for rows4.Next() {

				var music string
				err = rows4.Scan(&music)
				if err != nil {
					panic(err)
				}

				yt.Music = append(yt.Music, music)
			}

			//Assign to the post
			p.YouTubeData = *yt
		}

		posts[p.Id] = *p

	}

	//toJson(posts)
	//
	//data := readYAMLFile(postsFle)
	//posts := convertYAML(data)

	return posts
}

//func getTweets(tweetsFile string) map[string]Tweet {
//
//	data := readYAMLFile(tweetsFile)
//	tweets := convertTweetsYAML(data)
//
//	return tweets
//}

//func readYAMLFile(filename string) []byte {
//
//	data, err := ioutil.ReadFile(filename)
//
//	if err != nil {
//		log.Fatalf("Failed to read YML file : %v", err.Error())
//	}
//
//	return data
//}

type YouTubeData struct {
	Id    string
	Body  string
	Music []string
}

type Post struct {
	Id int
	Slug string
	Title       string
	Description string
	Date        string
	TopResult   string
	Keywords    []string
	YouTubeData YouTubeData
	Body        string
	Transcript  string
	ClickToTweet string
}
//
//
//func convertYAML(input []byte) map[string]Post {
//	posts := make(map[string]Post)
//
//	err := yaml.Unmarshal(input, &posts)
//	if err != nil {
//		log.Fatalf("error: %v", err)
//	}
//	return posts
//}
//
//type Tweet struct {
//	Link string
//}
//
//func convertTweetsYAML(input []byte) map[string]Tweet {
//	tweets := make(map[string]Tweet)
//
//	err := yaml.Unmarshal(input, &tweets)
//	if err != nil {
//		log.Fatalf("error: %v", err)
//	}
//	return tweets
//}
