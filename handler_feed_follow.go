package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lbergold-dev/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("usage: follow <url>")
	}

	ctx := context.Background()

	feed, err := s.db.GetFeedByUrl(ctx, cmd.arguments[0])
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("feed %q does not exist", cmd.arguments[0])
	}
	if err != nil {
		return fmt.Errorf("get feed: %w", err)
	}

	feedFollowsParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(ctx, feedFollowsParams)
	if err != nil {
		return fmt.Errorf("create feed follow: %w", err)
	}

	fmt.Printf("%s now follows feed %s\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	feedFollows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("get feed follow: %w", err)
	}

	if len(feedFollows) == 0 {
		fmt.Println("No feed follows found")
		return nil
	}

	for _, feedFollow := range feedFollows {
		fmt.Printf("* %s\n", feedFollow.FeedName)
	}

	return nil
}
