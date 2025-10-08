package main

import (
	"github.com/eyalhh/reddit-scraper/scraper"
)



func main() {
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
				Sha256Hash: "",
			},
		},
	}
	scraper.GetPosts(data, "")
}
