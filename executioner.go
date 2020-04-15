package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func execInput(input string) (bool, error) {
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")

	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return false, errors.New("Path is required")
		}
		err := os.Chdir(args[1])
		if err != nil {
			return false, err
		}
		return true, nil
	case "ls":
		currentDir, err := os.Getwd()

		if err != nil {
			return false, errors.New(err.Error())
		}
		files, err := ioutil.ReadDir(currentDir)

		if err != nil {
			return false, errors.New(err.Error())
		}

		newFiles := mapArray(files, func(item interface{}, i int) interface{} {
			file, ok := item.(os.FileInfo)

			if !ok {
				return errors.New(err.Error())
			}

			return file.Name()
		})

		returnString := ""

		for _, file := range newFiles {
			returnString += fmt.Sprintf("%v\n", file)
		}

		fmt.Println(returnString)
		return true, nil
	case "exit":
		os.Exit(0)
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()

	if err != nil {
		return false, err
	}

	return true, err
}
