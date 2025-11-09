package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Shell struct {
	Path   string
	Reader *bufio.Reader
}

func main() {
	// read the PATH environment variable
	path := os.Getenv("PATH")

	sh := &Shell{
		Path:   path,
		Reader: bufio.NewReader(os.Stdin),
	}

	sh.Run()
}

func (sh *Shell) Run() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		line, err := sh.Reader.ReadString('\n')
		if err != nil {
			break
		}

		// ignore empty input lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		command, args := sh.parseLine(line)

		switch command {
		case "exit":
			os.Exit(0)
		case "echo":
			sh.handleEcho(args)
		case "type":
			// ensure args[0] is safe to reference
			if len(args) == 0 {
				args = append(args, "")
			}
			sh.handleType(args, true)
		default:
			// ensure args[0] is safe to reference
			if len(args) == 0 {
				args = append(args, "")
			}
			args = append([]string{command}, args...)
			sh.handleType(args, false)
		}
	}
}

func (sh *Shell) parseLine(line string) (string, []string) {
	// trim newline and split by whitespace
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
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

func (sh *Shell) handleEcho(args []string) {
	// Print a newline-terminated line (tests expect lines).
	fmt.Fprintln(os.Stdout, strings.Join(args, " "))
}

func (sh *Shell) handleType(args []string, onlyCheck bool) {
	if len(args) > 0 {
		args[0] = strings.TrimSpace(args[0])
		if args[0] == "echo" || args[0] == "type" || args[0] == "exit" {
			fmt.Fprint(os.Stdout, args[0]+" is a shell builtin\n")
			return
		}

		if fullpath, found := sh.findCommandInPath(args[0]); found {
			if onlyCheck {
				fmt.Fprint(os.Stdout, args[0]+" is "+fullpath+"\n")
				return
			}

			// execute the command using the discovered full path
			cmd := exec.Command(fullpath, args[1:]...)
			out, _ := cmd.CombinedOutput()

			// normalize and print the captured output
			outStr := sh.normalizeOutput(out)
			fmt.Fprint(os.Stdout, outStr)
		} else {
			fmt.Fprint(os.Stdout, args[0]+": not found\n")
		}
	} else {
		// args[0] is not safe here; just print a generic not found line
		fmt.Fprint(os.Stdout, ": not found\n")
	}
}

func (sh *Shell) normalizeOutput(out []byte) string {
	// 1) remove a stray space before "\n " if present
	// 2) remove leading blank lines
	// 3) ensure output ends with a single newline
	outStr := string(out)
	outStr = strings.ReplaceAll(outStr, "\n ", "\n")
	outStr = strings.TrimLeft(outStr, "\n")
	if !strings.HasSuffix(outStr, "\n") {
		outStr += "\n"
	}
	return outStr
}

func (sh *Shell) findCommandInPath(cmd string) (string, bool) {
	if sh.Path == "" || cmd == "" {
		return "", false
	}
	paths := strings.Split(sh.Path, ":")
	for _, p := range paths {
		if p == "" {
			continue
		}
		fullPath := p + "/" + cmd
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}
		// ensure it's not a directory and is executable by someone
		if !info.IsDir() && info.Mode().Perm()&0111 != 0 {
			return fullPath, true
		}
	}
	return "", false
}

func (sh *Shell) handleDefault(line string) {
	// Remove the newline character
	line = line[:len(line)-1]
	// Simulate command not found
	fmt.Fprint(os.Stdout, line+": "+"command not found\n")
}
