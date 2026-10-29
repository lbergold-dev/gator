package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lbergold-dev/gator/internal/database"
)

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

func (f *RSSFeed) unescape() {
	f.Channel.Title = html.UnescapeString(f.Channel.Title)
	f.Channel.Description = html.UnescapeString(f.Channel.Description)

	for i, v := range f.Channel.Item {
		f.Channel.Item[i].Title = html.UnescapeString(v.Title)
		f.Channel.Item[i].Description = html.UnescapeString(v.Description)
	}
}

func scrapeFeeds(s *state, ctx context.Context) error {
	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return errors.New("no feeds followed")
	}

	if err != nil {
		return fmt.Errorf("get next feed to fetch: %w", err)
	}

	now := time.Now().UTC()
	err = s.db.MarkFeedFetchedById(ctx, database.MarkFeedFetchedByIdParams{
		UpdatedAt:     now,
		LastFetchedAt: sql.NullTime{Time: now, Valid: true},
		ID:            nextFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("mark feed as fetched: %w", err)
	}

	feed, err := fetchFeed(s, ctx, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("fetch feed: %w", err)
	}

	layouts := []string{
		time.RFC1123Z, // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC1123,  // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC822Z,  // "02 Jan 06 15:04 -0700"
		time.RFC822,   // "02 Jan 06 15:04 MST"
	}

	for _, item := range feed.Channel.Item {

		var pubDate time.Time
		sqlPubDate := sql.NullTime{Valid: false}

		for _, layout := range layouts {
			pubDate, err = time.Parse(layout, item.PubDate)
			if err == nil {
				break
			}
		}

		if err == nil {
			sqlPubDate = sql.NullTime{Time: pubDate, Valid: true}
		}

		description := sql.NullString{Valid: false}
		if item.Description != "" {
			description = sql.NullString{String: item.Description, Valid: true}
		}

		if item.Title == "" {
			item.Title = "(No title)"
		}

		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: sqlPubDate,
			FeedID:      nextFeed.ID,
		}

		if err = s.db.CreatePost(ctx, postParams); err != nil {
			log.Printf("create post from feed %s - item %s: %v", nextFeed.Url, item.Title, err)
		}
	}

	return nil
}

func fetchFeed(s *state, ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response error: %d", resp.StatusCode)
	}

	feed := &RSSFeed{}

	if err := xml.NewDecoder(resp.Body).Decode(feed); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	feed.unescape()

	return feed, nil
}
