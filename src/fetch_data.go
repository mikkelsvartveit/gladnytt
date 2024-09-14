package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title        string       `xml:"title"`
	Description  string       `xml:"description"`
	Link         string       `xml:"link"`
	PubDate      string       `xml:"pubDate"`
	MediaContent MediaContent `xml:"content"`
}

type MediaContent struct {
	URL    string `xml:"url,attr"`
	Type   string `xml:"type,attr"`
	Medium string `xml:"medium,attr"`
}

func fetchData() {
	url := "https://www.nrk.no/toppsaker.rss"

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var rss RSS
	err = xml.Unmarshal(body, &rss)

	if err != nil {
		fmt.Println("Error unmarshalling XML:", err)
		return
	}

	// Only fetch the first N items
	// rss.Channel.Items = rss.Channel.Items[:10]

	for _, item := range rss.Channel.Items {
		addArticleIfNotExists(item.Title, item.Description, item.Link, item.MediaContent.URL, item.PubDate)
	}
}
