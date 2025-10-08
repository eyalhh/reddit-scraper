package scraper

import (
	"bytes"
	"encoding/json"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

const (
	BASE_URL = "https://gql-fed.reddit.com/"
)

type RequestData struct {
	OperationName string            `json:"operationName"`
	Variables     RequestVariables  `json:"variables"`
	Extensions    RequestExtensions `json:"extensions"`
}
type RequestVariables struct {
	SubredditName            string `json:"subredditName"`
	Sort                     string `json:"sort"`
	Range                    string `json:"range"`
	After                    string `json:"after"`
	IncludeSubredditInPosts  bool   `json:"includeSubredditInPosts"`
	IncludePostStats         bool   `json:"includePostStats"`
	IncludeCurrentUserAwards bool   `json:"includeCurrentUserAwards"`
}

type RequestExtensions struct {
	PersistedQuery PersistedQuery `json:"persistedQuery"`
}

type PersistedQuery struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
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
	CommentCount int        `json:"commentCount"`
	AuthorInfo   AuthorInfo `json:"authorInfo"`
}

type AuthorInfo struct {
	Name string `json:"name"`
}

type Post struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	AuthorName   string `json:"authorName"`
	CommentCount int    `json:"commentCount"`
	CreatedAt    string `json:"createdAt"`
}

func GetPostsAtLength(data RequestData, token string, length int) ([]Post, error) {

	var currentLen int
	var res []Post

	for currentLen < length {
		posts, err := GetPosts(data, token, "")
		if err != nil {
			return nil, err
		}
		if len(posts) == 0 {
			return res, fmt.Errorf("not enough posts: %d/%d", len(res), length)
		}
		data.Variables.After = base64.StdEncoding.EncodeToString([]byte(posts[len(posts)-1].ID))
		res = append(res, posts...)
		currentLen += len(posts)
	}

	return res, nil

}

func GetPosts(data RequestData, token string, afterId string) ([]Post, error) {

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
