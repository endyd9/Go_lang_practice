package main

import (
	"errors"
	"fmt"
	"net/http"
)


type request struct{
	url string
	status string
}

var errRequestFaild = errors.New("Request Faild")

func main() {
	results := make(map[string]string)
	c := make(chan request)
	urls := []string{
		"https://www.airbnb.com/",
		"https://www.naver.com/",
		"https://www.amazon.com/",
		"https://www.reddit.com/",
		"https://www.google.com/",
		"https://soundcloud.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
		"https://academy.nomadcoders.co/",
	}
	for _, url := range urls {
		go hitURL(url, c)
	}

	for range urls {
		result := <-c
		results[result.url] = result.status
	}

	for url, status := range results {
		fmt.Println(url, status)
	}
	
}

func hitURL(url string, c chan<- request) {
	resp, err:= http.Get(url)
	status := "OK"
	if err != nil || resp.StatusCode >= 400 {
		status = "FAILED"
	}
	c <- request{url:url, status: status}
}