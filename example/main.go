package main

import (
	"fmt"
	m "github.com/keighl/metabolize"
	"net/http"
	"net/url"
)

type MetaData struct {
	Title       string  `meta:"og:title"`
	Description string  `meta:"og:description,description"`
	Type        string  `meta:"og:type"`
	URL         url.URL `meta:"og:url"`
	VideoWidth  int64   `meta:"og:video:width"`
	VideoHeight int64   `meta:"og:video:height"`
}

func main() {
	res, _ := http.Get("https://www.youtube.com/watch?v=FzRH3iTQPrk")

	data := new(MetaData)

	err := m.Metabolize(res.Body, data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Title: %s\n", data.Title)
	fmt.Printf("Description: %s\n", data.Description)
	fmt.Printf("Type: %s\n", data.Type)
	fmt.Printf("URL: %s\n", data.URL.String())
	fmt.Printf("VideoWidth: %d\n", data.VideoWidth)
	fmt.Printf("VideoHeight: %d\n", data.VideoHeight)
}
