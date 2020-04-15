package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"regexp"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		currentDir, err := os.Getwd()

		if err != nil {
			log.Fatalln(err)
		}

		hostname, err := os.Hostname()

		if err != nil {
			log.Fatalln(err)
		}

		user, err := user.Current()

		if err != nil {
			log.Fatalln(err)
		}

		nameRegexp := regexp.MustCompile("[ ]")

		userName := nameRegexp.ReplaceAllString(user.Name, "")

		fmt.Print(hostname + "@" + userName + " " + currentDir + "> ")

		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if ok, err := execInput(input); !ok {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
