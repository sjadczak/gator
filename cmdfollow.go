package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sjadczak/gator/internal/database"
)

func handleFollow(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		msg := " gator> Usage:\n" +
			" gator follow <url>\n" +
			" example: gator follow https://hnrss.org/newest"
		fmt.Println(msg)
		return ErrInvalidArgs
	}

	feedUrl := cmd.args[0]

	ctx := context.Background()
	user, err := s.db.GetUser(ctx, s.cfg.Username)
	if err != nil {
		msg := fmt.Sprintf(" gator> user `%s` not found\n", s.cfg.Username)
		return errors.New(msg)
	}
	feed, err := s.db.GetFeed(ctx, feedUrl)
	if err != nil {
		msg := fmt.Sprintf(" gator> no feed found for `%s`", feedUrl)
		return errors.New(msg)
	}
	follow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		msg := fmt.Sprintf(" gator> could not create follow: %s", err)
		return errors.New(msg)
	}

	fmt.Printf(" gator> `%s` is following `%s`\n", follow.Username, follow.Feedname)
	return nil
}

func handleFollowing(s *state, cmd command) error {
	ctx := context.Background()
	fmt.Printf(" gator> finding %s's feed follows...\n", s.cfg.Username)
	follows, err := s.db.GetUserFeedFollows(ctx, s.cfg.Username)
	if err != nil {
		msg := fmt.Sprintf(" gator> could not retrieve %s's follows", s.cfg.Username)
		return errors.New(msg)
	}

	if len(follows) == 0 {
		fmt.Printf(" gator> no follows found for `%s`\n", s.cfg.Username)
	} else {
		for _, f := range follows {
			fmt.Printf(" - %s: %s\n", f.Username, f.Feedname)
		}
	}
	return nil
}
