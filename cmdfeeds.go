package main

import (
	"context"
	"errors"
	"fmt"
)

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
