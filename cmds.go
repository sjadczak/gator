package main

import (
	"errors"

	"github.com/sjadczak/gator/internal/config"
	"github.com/sjadczak/gator/internal/database"
)

var (
	ErrUnknownCmd  = errors.New(" gator> unrecognized command")
	ErrInvalidArgs = errors.New(" gator> invalid command arguments")
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if fn, ok := c.cmds[cmd.name]; ok {
		return fn(s, cmd)
	}
	return ErrUnknownCmd
}

func (c *commands) register(name string, fn func(*state, command) error) error {
	c.cmds[name] = fn
	return nil
}
