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

	for {
		fmt.Fprint(os.Stdout, "$ ")
		reader := bufio.NewReader(os.Stdin)

		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		command := strings.Split(line, " ")[0]
		args := strings.Split(line, " ")[1:]
		switch command {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Fprint(os.Stdout, strings.Join(args, " "))
		case "type":
			if len(args) > 0 {
				args[0] = strings.TrimSpace(args[0])
				if args[0] == "echo" || args[0] == "type" || args[0] == "exit" {
					fmt.Fprint(os.Stdout, args[0]+" is a shell builtin\n")
				} else {
					fmt.Fprint(os.Stdout, args[0]+": not found\n")
				}
			} else {
				fmt.Fprint(os.Stdout, args[0]+": "+" not found\n")
			}
		default:
			// Remove the newline character
			line = line[:len(line)-1]
			// Simulate command not found
			fmt.Fprint(os.Stdout, line+": "+"command not found\n")
		}

	}
}
