package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(s, context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("fetch feed: %w", err)
	}

	fmt.Printf("%+v\n", feed)

	return nil
}
