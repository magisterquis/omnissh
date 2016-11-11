package main

/*
 * session.go
 * Handle session channels
 * By J. Stuart McMurray
 * Created 20161108
 * Last Modifed 20161110
 */

import (
	"encoding/binary"
	"fmt"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	SHELL       = "shell"
	EXEC        = "exec"
	EMPTYSTRING = ""
)

/* Session handles a session channel */
func Session(c ssh.Conn, sess ssh.Channel, reqs <-chan *ssh.Request) {
	defer sess.Close()

	/* Handle session requests as they come */
	var (
		ctypech = make(chan *ssh.Request)
		ctype   *ssh.Request
		term    *terminal.Terminal
		command string
		exit    uint32
		err     error
	)
	go handleSessionRequests(c, sess, reqs, ctypech, &term)

	/* Wait for a request for an exec or shell */
	ctype = <-ctypech
	switch ctype.Type {
	case "shell":
		Debug("[%v] Spawning shell", c.RemoteAddr())
		command = "shell"
		exit = HandleShell(&term, sess)
	case "exec":
		/* Pull out command to execute */
		command, _, err = stringFromPayload(ctype.Payload)
		if nil != err {
			m := fmt.Sprintf("Unable to read command: %v", err)
			Debug("[%v] %v", c.RemoteAddr(), m)
			fmt.Fprintf(sess, "%v", m)
			return
		}
		/* TODO: Run it, pipe in stdin/out/err */
		Debug("[%v] Executing command %q", c.RemoteAddr(), command)
		fmt.Fprintf(sess, "Executing %q\r\n", command) /* DEBUG */
	default:
		Debug("wtf?") /* DEBUG */

	}

	/* Return value to client */
	es := make([]byte, 4)
	binary.BigEndian.PutUint32(es, exit)
	sess.SendRequest("exit-status", false, es)
	Debug(
		"[%v] %v exited with status %v",
		c.RemoteAddr(),
		command,
		exit,
	)
}

/* handleSessionRequests handles the requests for a session */
func handleSessionRequests(
	c ssh.Conn,
	sess ssh.Channel,
	reqs <-chan *ssh.Request,
	ctypech chan<- *ssh.Request,
	term **terminal.Terminal,
) {
	var (
		err   error
		shell *string = &SHELL
		exec  *string = &EXEC
	)
	for req := range reqs {
		switch req.Type { /* TODO: Finish this */
		case "pty-req": /* Make a terminal */
			Debug(
				"[%v] PTY Request (%v / %q)",
				c.RemoteAddr(),
				req.WantReply,
				req.Payload,
			)
			if *term, err = newTerminal(
				sess,
				req.Payload,
			); nil != err {
				m := fmt.Sprintf(
					"Unable to create PTY: %v",
					err,
				)
				Debug("[%v] %v", c.RemoteAddr(), m)
				fmt.Fprintf(sess, "%v", m)
				req.Reply(false, []byte(m))
				continue
			}
			req.Reply(true, nil)
		case *shell: /* Spawn a shell */
			shell = &EMPTYSTRING
			req.Reply(true, nil)
			ctypech <- req
		case *exec: /* Execute a command */
			exec = &EMPTYSTRING
			req.Reply(true, nil)
			ctypech <- req
		default:
			Debug(
				"[%v] Unhandled session request %v %v %q",
				c.RemoteAddr(),
				req.Type,
				req.WantReply,
				req.Payload,
			)
			/* TODO: Have requests update terminal */
			req.Reply(false, nil)
		}

	}
}

/* newTerminal wraps the channel in a terminal */
func newTerminal(
	sess ssh.Channel,
	payload []byte,
) (*terminal.Terminal, error) {
	var (
		width  uint32
		height uint32
		err    error
	)
	/* New terminal */
	term := terminal.NewTerminal(sess, Prompt())
	/* Ignore terminal type */
	_, payload, err = stringFromPayload(payload)
	if nil != err {
		return nil, err
	}
	/* Terminal sizes in characters */
	width, payload, err = u32FromPayload(payload)
	if nil != err {
		return nil, err
	}
	height, payload, err = u32FromPayload(payload)
	if nil != err {
		return nil, err
	}
	/* Make the terminal, set its size */
	if 0 != width && 0 != height {
		term.SetSize(int(width), int(height))
	}
	return term, nil
}

/* u32FromPayload reads a uint32 from the given request payload.  It returns
 * the read uint32 as well as the remainder of the payload. */
func u32FromPayload(payload []byte) (uint32, []byte, error) {
	/* Make sure there's at least 32 bits */
	if 4 > len(payload) {
		return 0, nil, fmt.Errorf(
			"payload too small to read uint32 (%v < 4)",
			len(payload),
		)
	}
	return binary.BigEndian.Uint32(payload[:4]), payload[4:], nil
}

/* stringFromPayload reads a string length and a string from the given request
 * payload.  It returns the string as well as the remainder of the payload. */
func stringFromPayload(payload []byte) (string, []byte, error) {
	l, p, err := u32FromPayload(payload)
	if nil != err {
		return "", nil, err
	}
	/* Make sure payload is long enough */
	if uint32(len(p)) < l {
		return "", nil, fmt.Errorf(
			"payload too small to read string (%v < %v)",
			len(payload),
			l,
		)
	}
	return string(p[:l]), p[l:], nil
}
