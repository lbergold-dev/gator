package main

import (
	"log"
	"os"

	"github.com/lbergold-dev/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	s := &state{
		cfg: cfg,
	}

	input := os.Args[1:]
	if len(input) == 0 {
		log.Fatal("not enough arguments provided")
	}

	cmds := initCommands()

	cmd := command{
		name:      input[0],
		arguments: input[1:],
	}

	if err := cmds.run(s, cmd); err != nil {
		log.Fatal(err)
	}

}

func initCommands() commands {
	cmds := commands{
		handlerFuncs: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogins)

	return cmds
}
