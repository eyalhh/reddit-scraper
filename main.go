package main

import (
	"encoding/json"
	"log"
	"os"
	"fmt"
	"github.com/eyalhh/reddit-scraper/scraper"
	"github.com/joho/godotenv"
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

	authKey := os.Getenv("AUTH_TOKEN")
	hash := os.Getenv("HASH")
	data := scraper.RequestData{
		OperationName: "SubredditFeedElements",
		Variables: scraper.RequestVariables{
			SubredditName: subreddit,
			Sort: "HOT",
			Range: "ALL",
			IncludeSubredditInPosts: false,
			IncludePostStats: true,
			IncludeCurrentUserAwards: true,
		},
		Extensions: scraper.RequestExtensions{
			PersistedQuery: scraper.PersistedQuery{
				Version: 1,
				Sha256Hash: hash,
			},
		},
	}

	posts, err := scraper.GetPostsAtLength(data, authKey, numPosts)
	if err != nil {
		log.Fatal(err)
	}

	jsonPosts, err := json.MarshalIndent(posts, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("the posts are:\n" + string(jsonPosts))

}
