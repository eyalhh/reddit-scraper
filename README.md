# Initial Setup

I wanted to build the scraper in Go, as it’s the programming language I am most comfortable with for day-to-day coding.  
I set up HTTP Toolkit and the MEmu emulator on my Windows machine and downloaded Reddit’s latest version APK.  
I had a problem where I couldn’t log in to Reddit on the emulator because of reCAPTCHA. Thanks to Amit, I solved it by installing a lower version of the Reddit APK.  

---

## Reversing Reddit's API

After that, I began tampering with Reddit’s API and discovered a few important details by clicking some buttons on the screen:

1. **https://gql-fed.reddit.com/** is the base URL for fetching post data.  
2. The first request from Reddit to fetch post data of a given subreddit only yields about 25 posts or so.  
3. In the request for fetching the post data, there is a **sort** property that defines which posts to respond with and how to sort them.  

My main problem was fetching 100 posts as I was requested to do, and not 25.  
I wanted to know what the app does when 25 posts have been scrolled by the user and now it needs to load more. Then I discovered that it sends a special property in the request body called `after`.  

At first, it looked gibberish to me. After some thinking, I decided to Base64 decode it because of the equal signs at the end (Base64 padding).  
I got `t3_...`, which (remembering previous post data that I viewed) is the format for a Reddit post ID.  

I then figured that Reddit just called the same request to fetch new data but changed the `after` property to that of the last post ID that it already fetched (or more accurately, the Base64 encoding of that ID).  
I implemented the same functionality in my scraper, and it turned out great.  

When I tried various subreddits, I saw that for some of them, the post data, even after pulling, was not enough.  
I then remembered the sort property in the request body and decided to tamper with it. I used the filter option in the Reddit app, and it changed the sort property, causing a different post data collection to be queried.  

After some experiments, I discovered that if I set `sort: "HOT"`, it would result in a much, much bigger post collection than the default one (`sort: "TOP"`).  
I also decided to make my program work with different numbers of posts — not just to query 100 posts, but any given length (the program will throw an error if there aren’t enough posts in the HOT data collection).  

---

## The Program Itself

The program is essentially a Go module, with a `main.go` file at the root of the file tree and a `scraper` package that defines and implements the necessary scrape functionality.  

Go, being a compiled language, compiles to a self-contained static binary.  
The inputting of the subreddit name and the number of posts to query are just the first and second CLI arguments, respectively.  
