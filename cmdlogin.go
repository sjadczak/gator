package main

import (
	"context"
	"errors"
	"fmt"
)

func handleLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		msg := " gator> Usage:\n" +
			" gator login <username>\n" +
			" example: gator login boots\n"
		fmt.Println(msg)

		return ErrInvalidArgs
	}

	username := cmd.args[0]
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, username)
	if err != nil {
		msg := fmt.Sprintf(" gator> user `%s` not registered\n", username)
		return errors.New(msg)
	}
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf(" gator> '%s' logged in!\n", user.Name)
	return nil
}
