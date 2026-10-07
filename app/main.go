package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var builtIns map[string]func(args ...string)

func main() {
	builtIns = map[string]func(args ...string){
		"exit": func(args ...string) { os.Exit(0) },
		"echo": func(args ...string) { fmt.Println(strings.Join(args, " ")) },
		"pwd": func(args ...string) {
			workingDir, err := os.Getwd()

			if err != nil {
				fmt.Fprint(os.Stderr, "Error reading current dir", err)
				os.Exit(1)
			} else {
				fmt.Println(workingDir)
			}
		},
		"type": func(args ...string) {
			if len(args) == 0 {
				fmt.Println("type: missing operand")
				return
			}
			target := args[0]
			if _, exists := builtIns[target]; exists {
				fmt.Printf("%s is a shell builtin\n", target)
			} else {
				resultPath := findExecutable(target)
				if resultPath != "" {
					fmt.Println(target + " is " + resultPath)
				} else {
					fmt.Printf("%s: not found\n", target)
				}

			}

		},

		"cd": func(args ...string) {
			target := args[0]

			if strings.HasPrefix(target, "/") && len(target) != 0 {
				if dirExists(target) {
					os.Chdir(target)
				} else {

					fmt.Printf("cd: %s: No such file or directory\n", target)
				}

			} else {
				cleanPath := cleanFilepath(target)
				if dirExists(cleanPath) {
					os.Chdir(cleanPath)
				} else {

					fmt.Printf("cd: %s: No such file or directory\n", target)
				}
			}
		},
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")

		userInput, err := reader.ReadString('\n')

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

		handleCommand(command, arguments...)

	}

}

func findExecutable(target string) string {
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

func handleCommand(command string, args ...string) {
	if handler, exists := builtIns[command]; exists {
		handler(args...)
		return
	}

	if resultPath := findExecutable(command); resultPath != "" {
		cmd := exec.Command(command, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

	} else {
		fmt.Printf("%s: command not found \n", command)
	}

}

func dirExists(path string) bool {
	info, err := os.Stat(path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}

		return false
	}

	return info.IsDir()
}

func cleanFilepath(p string) string {
	cwd, _ := os.Getwd()

	for char := range strings.SplitSeq(p, "/") {

		if char == ".." {
			for i := len(cwd) - 1; i > 0; i-- {

				if string(cwd[i]) == "/" {

					cwd = cwd[:i]
					break
				}

			}
		} else if char == "." {
			continue

		} else {
			cwd += "/" + char
		}
	}

	return cwd

}
