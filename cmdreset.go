package main

import "context"

func handleReset(s *state, cmd command) error {
	ctx := context.Background()
	return s.db.Reset(ctx)
}
