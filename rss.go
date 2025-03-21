package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/avgra3/gator/internal/database"
	"html"
	"io"
	"log"
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

// Our aggregator
func scrapeFeeds(s *state, cmd command, user database.User) error {
	// Next feed to fetch
	ctx := context.Background()
	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}
	// Update the feeds table, with the upcoming feed id
	err = s.db.MarkFeedFetched(ctx, nextFeed.ID)
	if err != nil {
		return err
	}
	// Fetch the new feed
	//fetchFeed(ctx context.Context, feedURL string)
	rssFeeds, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return err
	}
	// Iterate over the items in a loop and print their titles to the console
	for _, rssFeed := range rssFeeds.Channel.Item {
		item := fmt.Sprintf("* %v", rssFeed.Title)
		log.Println(item)
	}
	return nil
}
