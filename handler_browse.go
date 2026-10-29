package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/lbergold-dev/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int

	if len(cmd.arguments) == 0 {
		limit = 2
	} else {
		var err error
		limit, err = strconv.Atoi(cmd.arguments[0])
		if err != nil {
			return fmt.Errorf("convert limit to int: %w", err)
		}

		if limit < 1 {
			return errors.New("limit must be greater than 0")
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{UserID: user.ID, Limit: int32(limit)})

	if err != nil {
		return fmt.Errorf("get posts for user: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("No posts found")
		return nil
	}

	for _, post := range posts {
		fmt.Println(post.Title + " - " + post.FeedName)
		if post.PublishedAt.Valid {
			fmt.Println(post.PublishedAt.Time)
		}
		if post.Description.Valid {
			fmt.Println(post.Description.String)
		}
		fmt.Println(post.Url)
		fmt.Println("--------------------------------")
	}

	return nil
}
