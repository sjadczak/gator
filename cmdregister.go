package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sjadczak/gator/internal/database"
)

func handleRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		msg := " gator> Usage:\n" +
			" gator register <username>\n" +
			" examples: gator register boots\n"
		fmt.Println(msg)
		return ErrInvalidArgs
	}

	username := cmd.args[0]
	ctx := context.Background()
	user, err := s.db.CreateUser(
		ctx,
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      username,
		},
	)
	if err != nil {
		msg := fmt.Sprintf(" gator> user `%s` already exists\n", username)
		return errors.New(msg)
	}
	s.cfg.SetUser(user.Name)
	fmt.Printf(" gator> user `%s` registered\n", user.Name)
	return nil
}
