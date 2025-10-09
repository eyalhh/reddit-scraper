package main

import (
	"encoding/json"
	"fmt"
	"github.com/eyalhh/reddit-scraper/scraper"
	"log"
	"os"
	"strconv"
)

const (
	HARDCODED_MOBILE_APP_LOID_SECRET="b2hYcG9xclpZdWIxa2c6"
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

	posts, err := scraper.GetPostsAtLength(data, HARDCODED_MOBILE_APP_LOID_SECRET, numPosts)
	if err != nil {
		log.Fatal(err)
	}

	jsonPosts, err := json.MarshalIndent(posts, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("the posts are:\n" + string(jsonPosts))

}
