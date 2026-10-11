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
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	userName := cmd.arguments[0]

	_, err := s.db.GetUser(context.Background(), userName)
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
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.arguments[0],
	}

	user, err := s.db.CreateUser(context.Background(), userParams)
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
	users, err := s.db.GetUsers(context.Background())
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
	if err := s.db.DeleteAllUser(context.Background()); err != nil {
		return fmt.Errorf("delete users: %w", err)
	}

	fmt.Println("All users got deleted!")

	return nil
}
