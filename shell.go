package main

/*
 * shell.go
 * Interactive shell
 * By J. Stuart McMurray
 * Created 20161110
 * Last Modified 20161110
 */

import (
	"bufio"
	"io"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

/* lineReader represents a way to read a line whether there's a terminal
allocated or not.  It also implements io.Writer. */
type lineReader struct {
	term   **terminal.Terminal
	sess   ssh.Channel
	reader bufio.Reader
}

/* newLineReader makes a lineRedare from the pointer to the terminal and the
shell session. */
func newlineReader(t **terminal.Terminal, s ssh.Channel) *lineReader {
	return &lineReader{
		term:   t,
		sess:   s,
		reader: *bufio.NewReader(s),
	}
}

/* ReadLine reads a line from l */
func (l *lineReader) ReadLine() (string, error) {
	/* Prefer the terminal */
	if nil != *l.term {
		(*l.term).SetPrompt(Prompt(*l.term))
		return (*l.term).ReadLine()
	}
	/* Failing that, read a line */
	line, err := l.reader.ReadString('\n')
	if nil != err {
		return line, err
	}
	/* Remove trailing crlf bytes */
	if '\n' == line[len(line)-1] {
		line = line[:len(line)-1]
	}
	if '\r' == line[len(line)-1] {
		line = line[:len(line)-1]
	}
	return line, err
}

/* Write wraps the underlying session's Write. */
func (l *lineReader) Write(buf []byte) (n int, err error) {
	n, err = l.sess.Write(buf)
	return
}

/* handleTerm handles a request for a shell */
func HandleShell(term **terminal.Terminal, sess ssh.Channel) uint32 {
	var (
		lr   = newlineReader(term, sess)
		line string /* Read line */
		exit bool
		err  error
	)

	/* TODO: Print welcome info */
	/* REPL */
	for {
		/* Grab a line from the user */
		line, err = lr.ReadLine()
		if nil != err {
			switch err {
			case io.EOF:
				return 0
			default:
				/* TODO: Catch errors better */
				Debug("Line read error: %v", err)
				return 1
			}
		}
		/* Run the command */
		_, _, exit = ExecuteCommand(line, sess, sess.Stderr())
		/* Exit if we ought */
		if exit {
			return 0
		}
	}
}
