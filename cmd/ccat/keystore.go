package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/secretprovider"
	"github.com/batmac/ccat/pkg/term"
)

func interactivelySetKey() {
	if !term.IsStdinTerminal() {
		log.Fatal("aborting, because stdin is not a terminal")
	}

	var name string
	var err error
	for len(name) == 0 {
		name, err = term.ReadLine("Key name: ")
		if err != nil {
			log.Fatal(err)
		}
	}

	secret, err := term.ReadPassword("Secret key: ")
	if err != nil {
		log.Fatal(err)
	}
	if len(secret) == 0 {
		log.Fatal("aborting, because no secret key was given")
	}

	fmt.Println("")

	log.Debugf("key name: %s, secret: %s\n", name, strings.Repeat("*", len(secret)))

	if value, err := secretprovider.GetSecret(name, ""); value != "" {
		log.Debugf("key %s already exists, overwritng...\n", name)
		log.Debugf("old secret: %s\n", strings.Repeat("*", len(value)))
	} else if err != nil && !errors.Is(err, secretprovider.ErrNotFound) {
		log.Fatal(err)
	}

	if err := secretprovider.SetSecret(name, secret); err != nil {
		log.Fatal(err)
	}

	log.Debugf("key successfully set")
	os.Exit(0)
}
