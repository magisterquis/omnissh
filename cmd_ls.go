package main

/*
 * cmd_ls.go
 * Lists a directory
 * By J. Stuart McMurray
 * Created 20161111
 * Last Modified 20161111
 */

import (
	"flag"
	"fmt"
	"io"
)

func init() {
	RegisterCommandFunc("ls", CommandLs)
}

/* CommandLs lists a directory */
func CommandLs(
	help bool,
	args []string,
	stdio io.ReadWriter,
	stderr io.ReadWriter,
) (usage string) {
	/* Obligatory help text */
	if help {
		return "ls [-h] [-l] [directory]"
	}
	/* Flags */
	fs := flag.NewFlagSet("ls", flag.ContinueOnError)
	long := fs.Bool("l", false, "Long-format listing")
	fs.Usage = func() {
		fmt.Fprintf(
			stderr,
			`Usage: ls [options] [directory]

Prints the contents of the given directory, or the current directory if no
directory is given.

Options:
`,
		)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); nil != err {
		Debug("PE %T: %v", err) /* DEBUG */
		return
	}
	fmt.Fprintf(stdio, "Listing %v (long:%v)", fs.Args(), long)
	return
}
