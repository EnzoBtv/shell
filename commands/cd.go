package commands

import (
	"errors"
	"os"
)

// Cd changes the current directory
func Cd(args []string) (bool, error) {
	if len(args) < 2 {
		return false, errors.New("Path is required")
	}
	err := os.Chdir(args[1])
	if err != nil {
		return false, err
	}
	return true, nil
}
