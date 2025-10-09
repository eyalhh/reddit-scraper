package scraper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/eyalhh/reddit-scraper/auth"
	"io"
	"net/http"
)

const (
	BASE_URL = "https://gql.reddit.com/"
)

type RequestData struct {
	ID        string           `json:"id"`
	Variables RequestVariables `json:"variables"`
}
type RequestVariables struct {
	SubredditName           string    `json:"subredditName"`
	Sort                    string    `json:"sort"`
	Range                   string    `json:"range"`
	After                   string    `json:"after"`
	AdContext               AdContext `json:"adContext"`
	IncludeSubredditInPosts bool      `json:"includeSubredditInPosts"`
	IncludePostStats        bool      `json:"includePostStats"`
}

type AdContext struct {
	Layout string `json:"layout"`
}

type ResponseData struct {
	Data InnerResponseData `json:"data"`
}

type InnerResponseData struct {
	PostFeed PostFeed `json:"postFeed"`
}

type PostFeed struct {
	Elements ResponseElements `json:"elements"`
}

type ResponseElements struct {
	Edges []ResponseEdges `json:"edges"`
}

type ResponseEdges struct {
	Node EdgeNode `json:"node"`
}

type EdgeNode struct {
	ID           string     `json:"id"`
	CreatedAt    string     `json:"createdAt"`
	Title        string     `json:"title"`
	CommentCount float64    `json:"commentCount"`
	AuthorInfo   AuthorInfo `json:"authorInfo"`
}

type AuthorInfo struct {
	Name string `json:"name"`
}

type Post struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	AuthorName   string  `json:"authorName"`
	CommentCount float64 `json:"commentCount"`
	CreatedAt    string  `json:"createdAt"`
}

func GetPostsAtLength(data RequestData, secret string, length int) ([]Post, error) {

	var currentLen int
	var res []Post
	var token string
	var err error

	for currentLen < length {
		token, err = auth.GetAccessToken(secret)
		if err != nil {
			return nil, err
		}
TryAgain:
		posts, err := GetPosts(data, token)
		if err != nil {
			if err.Error() == "invalid access_token" {
				fmt.Println("invalid access_token getting another one..")
				token, err = auth.GetAccessToken(secret)
				if err != nil {
					return nil, err
				}
				goto TryAgain
			}

			return nil, err
		}
		if len(posts) == 0 {
			return res, fmt.Errorf("not enough posts: %d/%d", len(res), length)
		}
		data.Variables.After = base64.StdEncoding.EncodeToString([]byte(posts[len(posts)-1].ID))
		res = append(res, posts...)
		currentLen += len(posts)
	}

	return res[:length], nil

}

func GetPosts(data RequestData, token string) ([]Post, error) {

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", BASE_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid access_token")
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseData ResponseData
	if err := json.Unmarshal(respBody, &responseData); err != nil {
		return nil, err
	}

	return formatPosts(responseData), nil
}

func formatPosts(responseData ResponseData) []Post {

	var posts []Post
	for _, post := range responseData.Data.PostFeed.Elements.Edges {
		var newPost Post
		newPost.ID = post.Node.ID
		newPost.CreatedAt = post.Node.CreatedAt
		newPost.Title = post.Node.Title
		newPost.CommentCount = post.Node.CommentCount
		newPost.AuthorName = post.Node.AuthorInfo.Name
		posts = append(posts, newPost)
	}

	return posts

}
