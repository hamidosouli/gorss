package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RssItem `xml:"item"`
	} `xml:"channel"`
}

type RssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (RSSFeed, error) { // used pointer so we can return RSSFeed{}
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}
	defer res.Body.Close()

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return RSSFeed{}, err
	}
	rssFeed := RSSFeed{}
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return RSSFeed{}, err
	}
	return rssFeed, nil
}
