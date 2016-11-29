package main

/*
 * cmd_cd.go
 * Change Directories
 * By J. Stuart McMurray
 * Created 20161111
 * Last Modified 20161111
 */

import (
	"fmt"
	"io"
	"os"
)

func init() {
	RegisterCommandFunc("cd", CommandCD)
}

/* CommandCD changes directories. */
func CommandCD(
	help bool,
	args []string,
	stdio io.ReadWriter,
	stderr io.ReadWriter,
) (usage string) {
	/* Obligatory help text */
	if help {
		return "cd <directory>"
	}
	/* Make sure we have a directory */
	if 0 == len(args) {
		fmt.Fprintf(
			stderr,
			"Usage: %v\r\n",
			CommandCD(true, nil, nil, nil),
		)
		return
	}
	/* Try */
	if err := os.Chdir(args[0]); nil != err {
		fmt.Fprintf(stderr, "%v\r\n", err)
	}
	return
}
