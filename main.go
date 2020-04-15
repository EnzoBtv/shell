package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		currentDir, err := os.Getwd()

		if err != nil {
			log.Fatalln(err)
		}

		fmt.Print(currentDir + "> ")

		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if ok, err := execInput(input); !ok {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
