package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/sjadczak/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("could not read config: %v", err)
		os.Exit(1)
	}

	s := &state{cfg}
	cmdMap := make(map[string]func(*state, command) error)
	cmds := &commands{
		cmdMap,
	}
	cmds.register("login", handleLogin)

	if len(os.Args) < 2 {
		msg := "missing command\n" +
			"Usage:\n" +
			"gator <command> [args...]"

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
		fmt.Printf("failed to run command: %v", err)
		os.Exit(1)
	}
}
