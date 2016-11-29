package main

/*
 * cmd_whoami.go
 * Prints the current user's info
 * By J. Stuart McMurray
 * Created 20161111
 * Last Modified 20161111
 */

import (
	"fmt"
	"io"
	"os/user"
)

func init() {
	RegisterCommandFunc("whoami", CommandWhoami)
}

/* CommandWhoami prints info about the user */
func CommandWhoami(
	help bool,
	args []string,
	stdio io.ReadWriter,
	stderr io.ReadWriter,
) (usage string) {
	/* Obligatory help text */
	if help {
		return "whoami"
	}
	/* Get and display current user */
	u, err := user.Current()
	if nil != err {
		fmt.Fprintf(stderr, "Unable to determine user info: %v", err)
		return
	}
	fmt.Fprintf(
		stdio,
		"Name: %v\r\n"+
			"User: %v\r\n"+
			"UID:  %v\r\n"+
			"Home: %v\r\n"+
			"GID:  %v\r\n",
		u.Name,
		u.Username,
		u.Uid,
		u.HomeDir,
		u.Gid,
	)
	return
}
