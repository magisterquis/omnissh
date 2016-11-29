package main

/*
 * config.go
 * Server config
 * By J. Stuart McMurray
 * Created 20161108
 * Last Modified 20161128
 */

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/subtle"
	"fmt"

	"golang.org/x/crypto/ssh"
)

const SERVERVERSION = "SSH-2.0-OpenSSH_7.0"
const AUTHKEY = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDORgHEnokyTo9DM4+NHX06PmQ/7AfMs80pEH+K8jPiJBVLv8urHJ/8CMVwFSCSssktkjf4im0aRw1SJX+QhJ0U3zmE02NPiepPGF0gqhCx5cD8Mdj4/ujETkKU9HsCGS/lEOp6rEkpd0GcohCndzQ4aSFW7MQ/Uy0GgCdWI3NsAQEwc9iJfJm74tlOv6S7hZBfZiFH0/kK26iv+rOPDHPPuF+akNNUXxzxi7h+MREhiqXNJjFigShXDWqjizvFofC3o9HUvSkvhojn/palBdLpB4c04MA7rOAnAJ4OGzqD3T8EyoKQ2T0UOHFOh3rtYo51jaaW5l+YhDXRGIEuPgzR stuart@stuart-mbp"

/* MakeServerConfig generates a config suitable to be used in SSH servers */
func MakeServerConfig() (*ssh.ServerConfig, error) {
	/* Function to check public keys */
	pkc, err := publicKeyChecker(AUTHKEY)
	if nil != err {
		return nil, err
	}

	/* Config to return */
	c := &ssh.ServerConfig{
		PublicKeyCallback: pkc,
		ServerVersion:     SERVERVERSION,
	}
	/* Generate an RSA key for the server */
	//rkey, err := rsa.GenerateKey(rand.Reader, 4096)
	rkey, err := rsa.GenerateKey(rand.Reader, 1024) /* DEBUG */
	if nil != err {
		return nil, err
	}
	/* Turn the RSA key into something usable by the server */
	skey, err := ssh.NewSignerFromKey(rkey)
	if nil != err {
		return nil, err
	}
	/* Print key in a usable form */
	Debug("Server's key: %s", ssh.FingerprintSHA256(skey.PublicKey()))
	/* Add the server's key */
	c.AddHostKey(skey)
	return c, nil
}

/* publicKeyChecker return a function which checks that key is acceptible */
func publicKeyChecker(authkey string) (func(
	ssh.ConnMetadata, ssh.PublicKey,
) (*ssh.Permissions, error), error) {
	/* Make sure the given authorized key parses */
	ak, _, _, _, err := ssh.ParseAuthorizedKey([]byte(authkey))
	if nil != err {
		return nil, err
	}
	/* Actual bytes to compare */
	bs := ak.Marshal()
	/* Return a function which checks it against the client's key */
	return func(
		conn ssh.ConnMetadata,
		key ssh.PublicKey,
	) (*ssh.Permissions, error) {
		/* Check keys */
		if 1 == subtle.ConstantTimeCompare(bs, key.Marshal()) {
			/* Return all is good if they match */
			return nil, nil
		}
		/* Tell the user if not */
		return nil, fmt.Errorf("no")
	}, nil
}
