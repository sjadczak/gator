package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/sjadczak/gator/internal/config"
	"github.com/sjadczak/gator/internal/database"
)

// main function runs the CLI
func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("could not read config: %v", err)
		os.Exit(1)
	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Printf("could not connect to db: %v", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	s := &state{dbQueries, cfg}
	cmdMap := make(map[string]func(*state, command) error)
	cmds := &commands{
		cmdMap,
	}
	cmds.register("login", handleLogin)
	cmds.register("register", handleRegister)
	cmds.register("reset", handleReset)
	cmds.register("users", handleUsers)
	cmds.register("agg", handleAgg)
	cmds.register("addfeed", handleAddFeed)
	cmds.register("feeds", handleFeeds)
	cmds.register("follow", handleFollow)
	cmds.register("following", handleFollowing)

	if len(os.Args) < 2 {
		msg := " gator> missing command\n\n" +
			" Usage:\n" +
			" gator <command> [args...]"

		fmt.Println(msg)
		os.Exit(1)
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(s, cmd)
	if errors.Is(err, ErrInvalidArgs) {
		os.Exit(1)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
