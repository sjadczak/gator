package main

import (
	"context"
	"errors"
	"fmt"
)

func handleUsers(s *state, cmd command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		msg := fmt.Sprintf(" gator> %v", err)
		return errors.New(msg)
	}

	fmt.Println(" gator> found users...")
	for _, user := range users {
		f := " - %s"
		if user.Name == s.cfg.Username {
			f += " (current)"
		}
		f += "\n"
		fmt.Printf(f, user.Name)
	}

	return nil
}
