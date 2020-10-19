package main

import (
	"github.com/viert/kstore/manager"
	"github.com/viert/kstore/term"
)

func main() {
	m := manager.New()
	err := m.Authenticate()
	if err != nil {
		term.Errorf("Error while authenticating: %s\n", err)
		return
	}
	m.Run()
}
