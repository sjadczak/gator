package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sjadczak/gator/internal/database"
)

// TODO: refactor to include feed follow on add new feed
func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		msg := " gator> Usage:\n" +
			" gator addfeed <name> <url>\n" +
			" example: gator addfeed \"Hacker News RSS\" \"https://hnrss.org/newest\""
		fmt.Println(msg)
		return ErrInvalidArgs
	}

	name := cmd.args[0]
	url := cmd.args[1]

	ctx := context.Background()
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
	_, err = s.db.CreateFeedFollow(
		ctx,
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return errors.New(" gator> failed to create feed follow")
	}

	fmt.Printf("%v\n", feed)
	return nil
}

func handleFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		msg := fmt.Sprintf(" gator> %v", err)
		return errors.New(msg)
	}

	fmt.Println(" gator> found feeds...")
	for _, feed := range feeds {
		fmt.Printf(" - %s, %s\n", feed.Name, feed.UserName)
	}
	return nil
}
