package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

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
				resultPath := findCommand(target)
				if resultPath != "" {
					fmt.Println(target + " is " + resultPath)
				} else {
					fmt.Printf("%s: not found\n", commandsAndArgs[1])
				}

			}
		default:
			resultPath := findCommand(command)

			if resultPath != "" {
				cmd := exec.Command(command, arguments...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()

			} else {
				fmt.Printf("%s: command not found \n", command)
			}

		}

	}

}

func findCommand(target string) string {
	pathEnv := os.Getenv("PATH")
	dirs := filepath.SplitList(pathEnv)

	for _, dir := range dirs {
		fullpath := filepath.Join(dir, target)
		info, err := os.Stat(fullpath)
		if err != nil {
			continue
		}

		if !info.IsDir() && info.Mode()&0111 != 0 {
			return fullpath
		}
	}
	return "" // If we get here, nothing was found
}
