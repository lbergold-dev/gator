package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/google/uuid"
	"github.com/lbergold-dev/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(s, context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("fetch feed: %w", err)
	}

	fmt.Printf("%+v\n", feed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arguments) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.name)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("create feed: %w", err)
	}

	fmt.Printf("created feed: %+v\n", feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("get feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', 0)

	for _, feed := range feeds {
		fmt.Fprintf(w, "* %s\t%s\t%s\n", feed.Name, feed.Url, feed.UserName)
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("flush writer: %w", err)
	}

	return nil
}
