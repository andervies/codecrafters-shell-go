package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type TokenType int

const (
	TokenError TokenType = iota
	TokenCommand
	TokenArgument
	// TokenPipe
	// TokenRedirect
)

type Token struct {
	Type  TokenType
	Value string
}

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

			if target == "~" {
				home, err := os.UserHomeDir()

				if err != nil {
					fmt.Println(err)
				}
				target = home

			} else {
				target = cleanFilepath(target)
			}

			if err := os.Chdir(target); err != nil {
				fmt.Printf("cd: %s: No such file or directory\n", target)

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

		tokens, err := Tokenize(userInput)
		if err != nil {
			fmt.Println(err)
			continue
		}

		command, arguments := ExtractCommands(tokens)

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

func cleanFilepath(p string) string {

	if strings.HasPrefix(p, "/") && len(p) != 0 {
		return p
	}
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

func Tokenize(input string) ([]Token, error) {

	var tokens []Token
	var buf strings.Builder

	inSingleQuotes := false
	inDoubleQuotes := false
	escapeNext := false
	// var quoteChar rune

	runes := []rune(input)

	for _, ch := range runes {

		if inSingleQuotes {
			if ch == '\'' {
				inSingleQuotes = false
			} else {
				buf.WriteRune(ch)
			}
		} else if inDoubleQuotes {
			if escapeNext {
				switch ch {
				case '"':
					buf.WriteRune(ch)
					escapeNext = false
				case '\\':
					buf.WriteRune(ch)
					escapeNext = false
				}
			} else if ch == '"' {
				inDoubleQuotes = false

			} else if ch == '\\' {
				escapeNext = true
			} else {
				buf.WriteRune(ch)
			}
		} else if escapeNext {
			buf.WriteRune(ch)
			escapeNext = false
		} else {
			switch ch {
			case '\'':
				inSingleQuotes = true

			case '"':
				inDoubleQuotes = true

			case ' ':
				if buf.Len() > 0 {
					tokens = append(tokens, createToken(buf.String(), tokens))
					buf.Reset()
				}

			case '\\':
				escapeNext = true

			// Add pipe case in future

			default:
				buf.WriteRune(ch)

			}
		}
	}

	if inSingleQuotes || inDoubleQuotes {

		return nil, fmt.Errorf("unclosed quote error")
	}

	if buf.Len() > 0 {
		tokens = append(tokens, createToken(buf.String(), tokens))
	}

	return tokens, nil
}

func createToken(value string, existingToken []Token) Token {

	// Add pipe check in the future
	if len(existingToken) == 0 {
		return Token{Type: TokenCommand, Value: value}

	}

	return Token{Type: TokenArgument, Value: value}
}

func ExtractCommands(tokens []Token) (string, []string) {
	var command string
	var args []string

	for _, token := range tokens {
		switch token.Type {
		case TokenCommand:
			command = token.Value

		case TokenArgument:
			args = append(args, token.Value)

		}
	}

	return command, args
}
