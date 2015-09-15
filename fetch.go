package main

import (
	"fmt"
	"net/http"
)

// Fetch makes a `GET` request to an
// array of URLs concurrently
func Fetch(urls []string) []*HTTPResponse {

	ch := make(chan *HTTPResponse)
	responses := []*HTTPResponse{}

	for _, url := range urls {
		go func(url string) {
			fmt.Printf("Fetching %s \n", url)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Error: ", err)
				ch <- &HTTPResponse{url, nil, err}
			} else {
				ch <- &HTTPResponse{url, resp, err}
			}
		}(url)
	}

	for {
		select {
		case r := <-ch:

			// Scrape
			Scrape(r)

			responses = append(responses, r)
			if len(responses) == len(urls) {
				return responses
			}

			defer r.response.Body.Close()
		}
	}

	return responses
}
