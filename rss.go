package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

// RSS Types
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// RSS functions
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Needs to fetch a feed from the givel URL
	// If nothing goes wrong, return a filled-out RSSFeed struct
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	// No we need to unmarshal the data
	feed := RSSFeed{}
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return &RSSFeed{}, err
	}

	// Need to run the Title and Description fields of both the entire channel and all the items
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	for _, rssFeed := range feed.Channel.Item {
		rssFeed.Title = html.UnescapeString(rssFeed.Title)
	}

	return &feed, nil
}
