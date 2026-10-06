package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	// TODO: Uncomment the code below to pass the first stage

	for {
		fmt.Print("$ ")
		text := bufio.NewReader(os.Stdin)
		userInput, err := text.ReadString('\n')

		if err != nil {
			fmt.Fprint(os.Stderr, "Error reading user input", err)
			os.Exit(1)
		}

		userInput = strings.TrimSpace(userInput)

		commandsAndArgs := strings.Fields(userInput)

		if len(commandsAndArgs) == 0 {
			continue
		}

		command := commandsAndArgs[0]

		arguments := commandsAndArgs[1:]

		switch command {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Print(strings.Join(arguments, " ") + "\n")
		case "type":
			if len(arguments) == 0 {
				fmt.Println("type: missing operand")
				continue
			}
			target := arguments[0]
			if target == "exit" || target == "echo" || target == "type" {
				fmt.Printf("%s is a shell builtin\n", commandsAndArgs[1])
			} else {
				pathEnv := os.Getenv("PATH")
				dirs := filepath.SplitList(pathEnv)
				found := false

				for _, dir := range dirs {

					fullpath := filepath.Join(dir, target)
					info, err := os.Stat(fullpath)

					if err != nil {

						continue
					}

					if !info.IsDir() && info.Mode()&0111 != 0 {
						fmt.Println(target + " is " + fullpath)
						found = true
						break
					}

				}

				if !found {
					fmt.Printf("%s: not found\n", commandsAndArgs[1])
				}

			}
		default:
			fmt.Printf("%s: command not found \n", command)
		}

	}

}
