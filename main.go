package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lbergold-dev/gator/internal/config"
	"github.com/lbergold-dev/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg        *config.Config
	db         *database.Queries
	httpClient *http.Client
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
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 20 * time.Second},
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
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	return cmds
}
