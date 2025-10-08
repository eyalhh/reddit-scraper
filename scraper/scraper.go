package scraper

import (
	"net/http"
	"encoding/json"
	"fmt"
	"bytes"
	"io"
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

func GetPosts(data RequestData, token string) error {

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonData))

	req, err := http.NewRequest("POST", BASE_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer " + token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Respnonse body:\n", string(respBody))


	return nil	
}

