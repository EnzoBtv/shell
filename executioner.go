package main

import (
	"os"
	"strings"

	"github.com/EnzoBtv/shell/commands"
)

func execInput(input string) (bool, error) {
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")

	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":
		return commands.Cd(args)
	case "ls":
		return commands.Ls()
	case "gitcontrib":
		return commands.Gitcontrib(args)
	case "exit":
		os.Exit(0)
	}

	// cmd := exec.Command(args[0], args[1:]...)

	// cmd.Stderr = os.Stderr
	// cmd.Stdout = os.Stdout

	// err := cmd.Run()

	// if err != nil {
	// 	return false, err
	// }

	return true, nil
}
