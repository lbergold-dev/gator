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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("usage: login <name>")
	}

	userName := cmd.arguments[0]

	ctx := context.Background()

	_, err := s.db.GetUserByName(ctx, userName)
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
		return errors.New("usage: register <name>")
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.arguments[0],
	}

	ctx := context.Background()

	user, err := s.db.CreateUser(ctx, userParams)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("set user: %w", err)
	}

	fmt.Printf("User %q has been created\n", user.Name)

	return nil
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()

	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("get users: %w", err)
	}

	for _, u := range users {
		if u.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", u.Name)
			continue
		}
		fmt.Printf("* %s\n", u.Name)
	}

	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()

	if err := s.db.DeleteAllUser(ctx); err != nil {
		return fmt.Errorf("delete users: %w", err)
	}

	if err := s.cfg.SetUser(""); err != nil {
		return fmt.Errorf("set user: %w", err)
	}

	fmt.Println("All users deleted")

	return nil
}
