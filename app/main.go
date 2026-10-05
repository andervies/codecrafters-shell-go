package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	// TODO: Uncomment the code below to pass the first stage

	for {
		fmt.Print("$ ")
		text := bufio.NewReader(os.Stdin)
		user_input, err := text.ReadString('\n')

		if err != nil {
			fmt.Fprint(os.Stderr, "Error reading user input", err)
			os.Exit(1)
		}

		user_input = strings.TrimSpace(user_input)

		commands_and_args := strings.Fields(user_input)

		command := commands_and_args[0]

		arguments := commands_and_args[1:]

		switch command {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Print(strings.Join(arguments, " ") + "\n")
		default:
			fmt.Printf("%s: command not found \n", command)
		}
	}

}
