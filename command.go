package main

/*
 * command.go
 * Register and manage built-in commands
 * By J. Stuart McMurray
 * Created 20161111
 * Last Modified 20161111
 */

import (
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/google/shlex"
)

/* Command represents a runnable built-in command.  If help is true, the
function should return a string with a short (one-line) synopsis.  Otherwise
the function should execute the command, treating stdio (which combines stdin
and stdout), and stderr as normal.  Any arguments passed to the command are
passed in in cmdline.  The returned string is ignored unless help is true. */
type CommandFunc func(
	help bool,
	args []string,
	stdio io.ReadWriter,
	stderr io.ReadWriter,
) (usage string)

/* commands stores the built-in commands by name */
var (
	CommandFuncs = map[string]CommandFunc{}
	commandslock = &sync.Mutex{}
)

/* RegisterCommandFunc registers a CommandFunc with the given name.  It panics
if there is already a function with that name.  It is meant to be called from
the init() function in a command function's source file. */
func RegisterCommandFunc(name string, f CommandFunc) {
	commandslock.Lock()
	defer commandslock.Unlock()
	if _, ok := CommandFuncs[name]; ok {
		log.Panicf("Command Function %q already registered", name)
	}
	CommandFuncs[name] = f
}

/* ExecuteCommand executes the command contained in line.  It returns whether
or not the session should exit and whether the command was found.  It takes the
commands stdio/err. */
func ExecuteCommand(
	line string,
	stdio io.ReadWriter,
	stderr io.ReadWriter,
) (unfound, unparse, exit bool) {
	/* Trim whitespace */
	line = strings.TrimSpace(line)
	/* Ignore comments and blank lines */
	if "" == line || '#' == line[0] {
		return
	}
	/* Split the command */
	parts, err := shlex.Split(line)
	if nil != err {
		fmt.Fprintf(stderr, "Unable to parse %q: %v", line, err)
		unparse = true
		return
	}
	Debug("Parts: %q", parts) /* DEBUG */
	/* Command is the first part */
	if 0 == len(parts) {
		return
	}
	command := parts[0]
	args := parts[1:]
	/* Special commands */
	switch command {
	case "exit", "quit", "bye":
		exit = true
		return
	}
	/* Get the function */
	cfunc, ok := CommandFuncs[command]
	if !ok {
		fmt.Fprintf(stdio, "Unknown command %q\r\n", command)
		unfound = true
		return
	}
	/* Run the function */
	cfunc(false, args, stdio, stderr)
	return
}
