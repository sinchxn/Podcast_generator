package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os" // To fetch environment variables

	twitterscraper "github.com/imperatrona/twitter-scraper"
)

type Tweet struct {
	Text string `json:"text"`
}

func main() {
	scraper := twitterscraper.New()

	// Get the token and CSRF token from environment variables
	token := os.Getenv("TWITTER_AUTH_TOKEN")
	csrfToken := os.Getenv("TWITTER_CSRF_TOKEN")

	if token == "" || csrfToken == "" {
		log.Fatal("Missing authentication tokens")
	}

	// Set up authentication using environment variables
	scraper.SetAuthToken(twitterscraper.AuthToken{
		Token:     token,
		CSRFToken: csrfToken,
	})

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
		for tweet := range scraper.SearchTweets(context.Background(), query, 2) {
			if tweet.Error != nil {
				log.Printf("Error: %v\n", tweet.Error)
				continue
			}
			tweets = append(tweets, Tweet{Text: tweet.Text})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tweets)
	})

	// Get the port from the environment variable (Render provides it dynamically)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback to 8080 if no port is provided
	}

	// Start the server
	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
