package main

import (
	"context"
	"encoding/xml"
	"errors"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/avgra3/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	for _, rssItem := range rssFeeds.Channel.Item {
		parsedDate, err := parsingDates(rssItem.PubDate)
		if err != nil {
			log.Println(err)
		}
		createPostParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       rssItem.Title,
			Url:         rssItem.Link,
			Description: rssItem.Description,
			PublishedAt: parsedDate,
			FeedID:      nextFeed.ID,
		}
		// If error where the post with the URL already exists, ignore it (can happen a lot)
		// If it's a different error, log it
		// Ensure the published_at field from the the feeds is correct, sometimes they're in different format than expected -- need to handle that
		// May need to manually convert the data into database/sql types
		_, err = s.db.CreatePost(ctx, createPostParams)
		ignoreError := &pq.Error{}
		if errors.As(err, &ignoreError) {
			// If we have an error due to duplicate url
			// Pass to next feed
			continue
		} else if err != nil {
			log.Println(createPostParams)
			log.Println(err)
			continue
		}
	}
	return nil
}

func parsingDates(dateToParse string) (time.Time, error) {
	// Pass 1
	layoutOne := "Jan 2, 2006 at 3:04pm (MST)"
	parsedDate, err := time.Parse(layoutOne, dateToParse)
	if err == nil {
		return parsedDate, nil
	}

	// Pass 2
	layoutTwo := "Mon, 02 Jan 2006 15:04:05 -0700"
	parsedDate, err = time.Parse(layoutTwo, dateToParse)
	if err == nil {
		return parsedDate, nil
	}

	// Pass 3
	layoutThree := "Thu, 20 Mar 2025 19:14:46 +000"
	parsedDate, err = time.Parse(layoutThree, dateToParse)
	if err == nil {
		return parsedDate, nil
	}

	// No luck :(
	log.Println("Unable to properly parse the given date :(")
	log.Printf("\nHERE IS THE TIME YOU WANTED TO PARSE:\n%v\n", dateToParse)
	return time.Time{}, nil
}
