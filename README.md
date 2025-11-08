go run ./app
echo is a shell builtin
# command-shell — Build your own shell in Go

This repository is a hands-on, step-by-step tutorial and implementation workspace for building a full-featured Unix-like shell in Go. The goal is to guide you through progressively implementing parsing, builtins, navigation, quoting, redirection, pipelines, autocompletion, persistent history, and more.

The code in `app/` provides starter code and small, focused examples you can extend. Use the roadmap below as your development checklist and reference while adding features.

## Goals

- Provide a minimal, easy-to-follow starting point (REPL, basic builtins).
- Implement a robust parser and execution engine that supports quoted arguments, redirection, and pipelines.
- Add common builtins and shell conveniences (cd, pwd, history, completion).
- Persist shell history and provide a pleasant interactive UX (arrow navigation, completions).
- Keep the codebase small, well-tested, and easy to read so it can be used as a learning resource.

## Files

- `app/main.go` — starter REPL and simple builtin examples
- `your_program.sh` — example shell script included for testing or demo purposes
- `go.mod` — Go module file

## Requirements

- Go 1.16+ (modern Go toolchain)
- A POSIX-like terminal (macOS, Linux, WSL on Windows)

## Build & run

From the project root (`/Users/arihant/Documents/command-shell`) you can run the program directly or build a binary.

Run without a build:

```bash
# run the program using go
cd /Users/arihant/Documents/command-shell
go run ./app
```

Build and run the binary:

```bash
cd /Users/arihant/Documents/command-shell
go build -o command-shell ./app
# then run
./command-shell
```

## Usage example (starter behavior)

The starter code in `app/main.go` implements a simple prompt and a few builtins. Example session:

```
$ echo Hello world
Hello world
$ type echo
echo is a shell builtin
$ type ls
ls: not found
$ foobar
foobar: command not found
$ exit
```

This repository's purpose is to evolve this starter shell into a complete implementation following the roadmap below.

## Roadmap / Feature checklist

Use this as your project plan. Each item is a small, testable milestone.

### Core shell basics

- [ ] Introduction
- [ ] Repository Setup
- [ ] Print a prompt
- [ ] Handle invalid commands
- [ ] Implement a REPL
- [ ] Implement `exit`
- [ ] Implement `echo`
- [ ] Implement `type`
- [ ] Locate executable files (PATH search)
- [ ] Run a program (spawn external processes)

### NAVIGATION

- [ ] The `pwd` builtin
- [ ] The `cd` builtin: Absolute paths
- [ ] The `cd` builtin: Relative paths
- [ ] The `cd` builtin: Home directory (`~`)

### QUOTING

- [ ] Single quotes: literal strings
- [ ] Double quotes: allow escapes and expansions
- [ ] Backslash outside quotes: escaping
- [ ] Backslash within single quotes: literal behavior
- [ ] Backslash within double quotes: escape handling
- [ ] Executing a quoted executable (paths with spaces)

### REDIRECTION

- [ ] Redirect stdout (`>`)
- [ ] Redirect stderr (`2>`)
- [ ] Append stdout (`>>`)
- [ ] Append stderr (`2>>`)

### AUTOCOMPLETION

- [ ] Builtin completion
- [ ] Completion with arguments
- [ ] Missing completions handling
- [ ] Executable completion from `$PATH`
- [ ] Multiple completions listing
- [ ] Partial completions (tab completion)

### PIPELINES

- [ ] Dual-command pipeline (`cmd1 | cmd2`)
- [ ] Pipelines with built-ins (where applicable)
- [ ] Multi-command pipelines (`cmd1 | cmd2 | cmd3`)

### HISTORY

- [ ] The `history` builtin
- [ ] Listing history
- [ ] Limiting history entries (size cap)
- [ ] Up-arrow navigation
- [ ] Down-arrow navigation
- [ ] Executing commands from history (`!n` style)

### HISTORY PERSISTENCE

- [ ] Read history from file on startup
- [ ] Write history to file on exit
- [ ] Append history to file (avoid truncation)
- [ ] Read history on startup (merge with session)
- [ ] Write history on exit (persist new entries)
- [ ] Append history on exit (safe concurrent appends)

## Contributing

This repository is organized as a tutorial: implement one checklist item at a time. For each feature:

- Add a small, focused test where practical (parsing, builtin behavior).
- Keep changes minimal and well-documented.
- Open PRs that implement a single checklist item.

Suggested starter tasks:

- Implement `pwd` and `cd` builtins.
- Improve the argument parser to respect quoting and escapes.
- Add tests for the `echo` and `type` builtins.

If you'd like, I can implement any checklist item for you and add tests and usage examples.

## License

Add a `LICENSE` file if you want to publish this project with a specific license (MIT, Apache-2.0, etc.).

