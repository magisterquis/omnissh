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
	"fmt"

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
		(*l.term).SetPrompt(Prompt())
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
	lr := newlineReader(term, sess)
	l, err := lr.ReadLine()
	if nil != err {
		Debug("E: %v", err)
		return 0
	}
	fmt.Fprintf(lr, "%v<--\r\n", l)
	fmt.Fprintf(lr, "Done.\r\n")
	return 0
}
