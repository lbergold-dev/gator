package main

import "fmt"

func handlerLogins(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("login command requires exactly 1 argument, got %d", len(cmd.arguments))
	}

	user := cmd.arguments[0]

	if err := s.cfg.SetUser(user); err != nil {
		return fmt.Errorf("set user: %w", err)
	}

	fmt.Printf("User has been set to %s\n", user)

	return nil
}
