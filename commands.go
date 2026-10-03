package main

import "fmt"

type commands struct {
	handlerFuncs map[string]func(*state, command) error
}

type command struct {
	name      string
	arguments []string
}

func (c *commands) run(s *state, cmd command) error {
	cmdHandler, ok := c.handlerFuncs[cmd.name]
	if !ok {
		return fmt.Errorf("command %s not found", cmd.name)
	}

	if err := cmdHandler(s, cmd); err != nil {
		return fmt.Errorf("run handler: %w", err)
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlerFuncs[name] = f
}
