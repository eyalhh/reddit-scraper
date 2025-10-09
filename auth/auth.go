package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const (
	BASE_URL = "https://accounts.reddit.com/"
)

type RequestData struct {
	Scopes []string `json:"scopes"`
}

type ResponseData struct {
	AccessToken string `json:"access_token"`
}

func GetAccessToken(secret string) (string, error) {

	reqData := RequestData{
		Scopes: []string{"*", "email", "pii"},
	}

	jsonData, err := json.MarshalIndent(reqData, "", "    ")
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", BASE_URL+"api/access_token", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Basic "+secret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}


	var responseData ResponseData
	if err := json.Unmarshal(respBody, &responseData); err != nil {
		return "", err
	}

	return responseData.AccessToken, nil

}
