package main

/*
 * prompt.go
 * Handles making the prompt for the user
 * By J. Stuart McMurray
 * Created 20161110
 * Last Modified 20161110
 */

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/crypto/ssh/terminal"
)

/* PromptData holds the cached prompt data */
var promptData struct {
	user string
	host string
	head string
	tail string
}

func init() {
	/* $ or #, as appropriate */
	uid := os.Getuid()
	if 0 == uid {
		/* TODO: Handle Windows */
		promptData.tail = "]# "
	} else {
		promptData.tail = "]$ "
	}
	/* Initial prompt username */
	u, err := user.Current()
	if nil == err {
		promptData.user = u.Username
	} else {
		Debug("Unable to get username for prompt: %v", err)
		promptData.user += fmt.Sprintf("%v", uid)
	}
	/* Initial prompt hostname */
	h, err := os.Hostname()
	if nil == err {
		promptData.host = h
	} else {
		Debug("Unable to get hostname for prompt: %v", err)
		promptData.host = "<unknown>"
	}
	regenPrompt()
}

/* Prompt returns a string like [user@hostname:directory]$ which is
valid for the current state of the process. */
func Prompt(t *terminal.Terminal) string {
	/* Get working directory */
	wd, err := os.Getwd()
	if nil != err {
		wd = err.Error()
	} else {
		wd = filepath.Clean(wd)
	}
	/* TODO: Colors */
	return string(t.Escape.Red) +
		promptData.head +
		wd +
		promptData.tail +
		string(t.Escape.Reset)
}

/* SetPromptUser sets the username for the prompt */
func SetPromptUser(u string) {
	promptData.user = u
	regenPrompt()
}

/* SetPromptHost sets the hostname for the prompt */
func SetPromptHost(h string) {
	promptData.host = h
	regenPrompt()
}

/* regenPrompt regenerates the [user@host: part of the prompt */
func regenPrompt() {
	promptData.head = "[" + promptData.user + "@" + promptData.host + ":"
}
