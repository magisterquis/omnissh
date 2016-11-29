package main

/*
 * cmd_help.go
 * Lists built-in commands
 * By J. Stuart McMurray
 * Created 201601111
 * Last Modified 201601111
 */

import (
	"fmt"
	"io"
	"sort"
)

func init() {
	RegisterCommandFunc("help", CommandHelp)
}

/* CommandHelp lists all the commands */
func CommandHelp(
	help bool,
	args []string,
	stdio io.ReadWriter,
	stderr io.ReadWriter,
) (usage string) {
	/* Obligatory help text */
	if help {
		return "help [command]"
	}
	/* If commands were given, print their help only */
	if 0 != len(args) {
		for _, arg := range args {
			f, ok := CommandFuncs[arg]
			if !ok {
				fmt.Fprintf(
					stderr,
					"Command %q unhelpfully unknown",
					arg,
				)
				continue
			}
			fmt.Fprintf(stdio, "%v\r\n", f(true, nil, nil, nil))
		}
		return
	}
	/* Get all the defined functions */
	ns := make([]string, 0, len(CommandFuncs))
	for k := range CommandFuncs {
		ns = append(ns, k)
	}
	sort.Strings(ns)
	/* Print help */
	fmt.Fprintf(stdio, "Defined commands:\r\n\r\n")
	for _, n := range ns {
		fmt.Fprintf(
			stdio,
			"%v\r\n",
			CommandFuncs[n](true, nil, nil, nil),
		)
	}
	return
}
