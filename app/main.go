package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var _ = fmt.Fprint
var _ = os.Stdout

func main() {

	// read the PATH environment variable
	path := os.Getenv("PATH")

	// start the REPL loop
	repl(path)
}

func repl(path string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "$ ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		command, args := parseLine(line)

		switch command {
		case "exit":
			os.Exit(0)
		case "echo":
			handleEcho(args)
		case "type":
			// ensure args[0] is safe to reference in the original fmt lines
			if len(args) == 0 {
				args = append(args, "")
			}
			handleType(args, path)
		default:
			handleDefault(line)
		}
	}
}

func parseLine(line string) (string, []string) {
	// mimic original splitting by space (retain behavior)
	parts := strings.Split(line, " ")
	if len(parts) == 0 {
		return "", nil
	}
	command := parts[0]
	args := []string{}
	if len(parts) > 1 {
		args = parts[1:]
	}
	return command, args
}

func handleEcho(args []string) {
	fmt.Fprint(os.Stdout, strings.Join(args, " "))
}

func handleType(args []string, path string) {
	if len(args) > 0 {
		args[0] = strings.TrimSpace(args[0])
		if args[0] == "echo" || args[0] == "type" || args[0] == "exit" {
			fmt.Fprint(os.Stdout, args[0]+" is a shell builtin\n")
		} else {
			// check if this command exists in the PATH (extracted to helper)
			if rel, found := findCommandInPath(path, args[0]); found {
				fmt.Fprintln(os.Stdout, args[0]+" is "+rel)
			} else {
				fmt.Fprint(os.Stdout, args[0]+": not found\n")
			}
		}
	} else {
		fmt.Fprint(os.Stdout, args[0]+": not found\n")
	}
}

// new helper: search PATH for a command and return its full path if found
func findCommandInPath(path, cmd string) (string, bool) {
	if path == "" || cmd == "" {
		return "", false
	}
	paths := strings.SplitSeq(path, ":")
	for p := range paths {
		if p == "" {
			continue
		}
		fullPath := p + "/" + cmd
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}
		// ensure it's not a directory and is executable by someone
		if !info.Mode().IsDir() && info.Mode()&0111 != 0 {
			return fullPath, true
		}
	}
	return "", false
}

func handleDefault(line string) {
	// Remove the newline character
	line = line[:len(line)-1]
	// Simulate command not found
	fmt.Fprint(os.Stdout, line+": "+"command not found\n")
}
