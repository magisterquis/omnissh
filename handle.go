package main

/*
 * handle.go
 * Handle SSH connections
 * By J. Stuart McMurray
 * Created 20161108
 * Last Modified 20161108
 */

import (
	"fmt"
	"io"
	"net"

	"golang.org/x/crypto/ssh"
)

/* Handle handles a client connection */
func Handle(c net.Conn, conf *ssh.ServerConfig) {
	/* Upgrade to an SSH connection */
	sc, chans, reqs, err := ssh.NewServerConn(c, conf)
	if nil != err {
		Debug(
			"[%v] Auth fail: %v",
			c.RemoteAddr(),
			err,
		)
		c.Close()
		return
	}
	/* Log when connection's closed */
	var cerr error
	defer func() {
		sc.Close()
		m := fmt.Sprintf(
			"[%v] Connection finished",
			sc.RemoteAddr(),
		)
		if nil != cerr && io.EOF != cerr {
			m += " (" + cerr.Error() + ")"
		}
		Debug("%s", m)
	}()
	Debug(
		"[%v] Authenticated as %v",
		sc.RemoteAddr(),
		sc.User(),
	)
	/* Print out channels and requests */
	go handleChans(sc, chans)
	go handleReqs(sc, reqs)
	/* Wait for connection to close, save the error for the logs */
	cerr = sc.Wait()
}

/* handleReqs handles requests for a client */
func handleReqs(c ssh.Conn, reqs <-chan *ssh.Request) {
	for req := range reqs {
		switch req.Type {
		case "no-more-sessions@openssh.com":
			Debug(
				"[%v] Replying no to request %v (%v / %q)",
				c.RemoteAddr(),
				req.Type,
				req.WantReply,
				req.Payload,
			)
			req.Reply(false, nil)
		default:
			Debug(
				"[%v] Unhandled request: %v (%v / %q)",
				c.RemoteAddr(),
				req.Type,
				req.WantReply,
				req.Payload,
			)
		}
	}
}

/* handleChans handles channels for a client */
func handleChans(c ssh.Conn, chans <-chan ssh.NewChannel) {
	for ch := range chans {
		switch ch.ChannelType() {
		case "session":
			Debug(
				"[%v] New session channel (%q)",
				c.RemoteAddr(),
				ch.ExtraData(),
			)
			sess, reqs, err := ch.Accept()
			if nil != err {
				ch.Reject(
					ssh.ResourceShortage,
					err.Error(),
				)
				continue
			}
			go Session(c, sess, reqs)
		default:
			Debug(
				"[%v] Rejecting %v channel request (%q)",
				c.RemoteAddr(),
				ch.ChannelType(),
				ch.ExtraData(),
			)
			ch.Reject(ssh.UnknownChannelType, "womp womp")
		}
	}
}
