package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sjadczak/gator/internal/database"
)

func handleAddFeed(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		msg := " gator> Usage:\n" +
			" gator addfeed <name> <url>\n" +
			"example: gator addfeed \"Hacker News RSS\" \"https://hnrss.org/newest\""
		fmt.Println(msg)
		return ErrInvalidArgs
	}

	if len(cmd.args) < 2 {
		return errors.New("gator addfeed requires 2 args, `name` and `feed url`")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	ctx := context.Background()
	user, err := s.db.GetUser(ctx, s.cfg.Username)
	if err != nil {
		msg := fmt.Sprintf(" gator> user `%s` not found\n")
		return errors.New(msg)
	}
	feed, err := s.db.CreateFeed(
		ctx,
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
			Url:       url,
			UserID:    user.ID,
		},
	)
	if err != nil {
		msg := fmt.Sprintf(" gator> failed to create feed `%s`", name)
		return errors.New(msg)
	}

	fmt.Printf("%v\n", feed)
	return nil
}
