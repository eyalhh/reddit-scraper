package main

import (
	"github.com/eyalhh/reddit-scraper/scraper"
	"log"
	"github.com/joho/godotenv"
	"os"
)



func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	authKey := os.Getenv("AUTH_TOKEN")
	hash := os.Getenv("HASH")
	data := scraper.RequestData{
		OperationName: "SubredditFeedElements",
		Variables: scraper.RequestVariables{
			SubredditName: "funny",
			Sort: "HOT",
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

	scraper.GetPosts(data, authKey)
}
