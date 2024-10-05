package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	twitterscraper "github.com/imperatrona/twitter-scraper"
)

type Tweet struct {
	Text string `json:"text"`
}

func main() {
	scraper := twitterscraper.New()

	// Set up authentication (you'll need to replace these with your actual credentials)
	scraper.SetAuthToken(twitterscraper.AuthToken{Token: "9e918c71af7b13dfbd8e24daeeca6ce655905653", CSRFToken: "ac60b58809c0705f0611997e4527696f8d9cca79988ff0e9c9440835f00de72026d051ad925fdd2a3c5c71405d747bf6ad3e59bc5fda9021709fae4b298c48371be7633aa251575c07d8f5aaf288c4ac"})

	// Check if login is successful
	if !scraper.IsLoggedIn() {
		log.Fatal("Failed to log in")
	}

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "Missing query parameter 'q'", http.StatusBadRequest)
			return
		}

		tweets := []Tweet{}
		for tweet := range scraper.SearchTweets(context.Background(), query, 100) {
			if tweet.Error != nil {
				log.Printf("Error: %v\n", tweet.Error)
				continue
			}
			tweets = append(tweets, Tweet{Text: tweet.Text})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tweets)
	})

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}