package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/google/uuid"
	"github.com/lbergold-dev/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("usage: agg <time_between_requests>")
	}

	timeBetweenRequest, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("parse duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)
	defer ticker.Stop()

	ctx := context.Background()

	for {
		if err := scrapeFeeds(s, ctx); err != nil {
			log.Println(err)
		}
		<-ticker.C
	}

}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.name)
	}

	ctx := context.Background()

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(ctx, feedParams)
	if err != nil {
		return fmt.Errorf("create feed: %w", err)
	}

	fmt.Printf("Created feed %s with url %q\n", feed.Name, feed.Url)

	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("create feed follow: %w", err)
	}

	fmt.Println("Followed feed")

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("get feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found")
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
