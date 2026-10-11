package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
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
