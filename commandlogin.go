package main

import "fmt"

func handleLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		msg := "Usage:\n" +
			"gator login <username>\n" +
			"example: gator login boots\n"
		fmt.Println(msg)

		return ErrInvalidArgs
	}

	username := cmd.args[0]
	err := s.cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("'%s' logged in!\n", username)
	return nil
}
