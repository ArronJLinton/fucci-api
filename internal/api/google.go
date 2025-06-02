package api

import (
	"fmt"
	"log"
	"net/http"

	customsearch "google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi/transport"
)

const (
	apiKey = "AIzaSyCX9JKYZ1-SeaQYa4GfHqUPfULWzHKHFdY"
	cx     = "some-custom-search-engine-id"
	query  = "arsenal"
)

func (c *Config) search(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{Transport: &transport.APIKey{Key: apiKey}}

	svc, err := customsearch.New(client)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := svc.Cse.List().Cx(cx).Q(query).Do()
	if err != nil {
		log.Fatal(err.Error())
	}

	for i, result := range resp.Items {
		fmt.Printf("#%d: %s\n", i+1, result.Title)
		fmt.Printf("\t%s\n", result.Snippet)
	}
	// respondWithJSON(w, http.StatusOK, result)
}
