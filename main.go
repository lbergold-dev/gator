package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lbergold-dev/gator/internal/config"
	"github.com/lbergold-dev/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	input := os.Args[1:]
	if len(input) == 0 {
		fmt.Fprintln(os.Stderr, "not enough arguments provided")
		os.Exit(1)
	}

	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	programState := &state{
		cfg: cfg,
	}

	db, err := sql.Open("postgres", programState.cfg.DBURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()
	programState.db = database.New(db)

	cmds := initCommands()
	cmd := command{
		name:      input[0],
		arguments: input[1:],
	}

	if err := cmds.run(programState, cmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func initCommands() commands {
	cmds := commands{
		handlerFuncs: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogins)
	cmds.register("register", handlerRegister)

	return cmds
}
