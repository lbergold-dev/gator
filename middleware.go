package main

import (
	"context"
	"fmt"

	"github.com/lbergold-dev/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		user, err := s.db.GetUserByName(ctx, s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("get user by name: %w", err)
		}

		return handler(s, cmd, user)
	}
}
