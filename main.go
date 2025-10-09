package main

import (
	"encoding/json"
	"fmt"
	"github.com/eyalhh/reddit-scraper/scraper"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <subreddit_name> <num_posts>\n", os.Args[0])
		return
	}
	subreddit := os.Args[1]
	numPosts, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("second argument must be integer\n")
		return
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secret := os.Getenv("MOBILE_APP_LOID_SECRET_BASE64")
	data := scraper.RequestData{
		ID: "5fb5c243c4c4",
		Variables: scraper.RequestVariables{
			SubredditName:           subreddit,
			Sort:                    "HOT",
			Range:                   "ALL",
			AdContext: scraper.AdContext{
				Layout: "CARD",
			},
			IncludeSubredditInPosts: false,
			IncludePostStats:        true,
		},
	}

	posts, err := scraper.GetPostsAtLength(data, secret, numPosts)
	if err != nil {
		log.Fatal(err)
	}

	jsonPosts, err := json.MarshalIndent(posts, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("the posts are:\n" + string(jsonPosts))

}
