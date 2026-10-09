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

func handlerLogins(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("login command requires exactly 1 argument, got %d", len(cmd.arguments))
	}

	userName := cmd.arguments[0]

	_, err := s.db.GetUser(context.Background(), sql.NullString{String: userName, Valid: true})
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("user %q does not exist", userName)
	}
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	if err := s.cfg.SetUser(userName); err != nil {
		return fmt.Errorf("set user: %w", err)
	}

	fmt.Printf("User has been set to %s\n", userName)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("register command requires exactly 1 argument, got %d", len(cmd.arguments))
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      sql.NullString{String: cmd.arguments[0], Valid: true},
	}

	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	if err := s.cfg.SetUser(user.Name.String); err != nil {
		return fmt.Errorf("set user: %w", err)
	}

	fmt.Printf("User %q has been created\n", user.Name.String)

	return nil
}
