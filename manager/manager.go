package manager

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"
	"github.com/viert/kstore/sec"
	"github.com/viert/kstore/store"
	"github.com/viert/kstore/term"
)

type cmdHandler func(string, string, ...string)

// Manager represents a console manager object
type Manager struct {
	rl    *readline.Instance
	enc   sec.Crypter
	store *store.Store

	accessToken string

	stopped  bool
	data     serviceData
	handlers map[string]cmdHandler
}

// New creates a new Manager instance
func New() *Manager {
	m := &Manager{
		stopped: false,
	}
	m.setupHandlers()
	m.setupReadline()
	return m
}

func (m *Manager) cmdLoop() {
	for !m.stopped {
		line, err := m.rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			m.stopped = true
		}
		m.cmd(line)
	}
}

func (m *Manager) setPrompt() {
	m.rl.SetPrompt(term.Blue("kstore") + "> ")
}

func (m *Manager) load() error {
	data, err := m.store.Load()
	if err != nil {
		if store.IsFileNotFound(err) {
			fmt.Println("no previous data found, initializing empty store...")
			m.data = make(serviceData)
			return nil
		}
		return err
	}

	dec, err := m.enc.Decrypt(data)
	if err != nil {
		return fmt.Errorf("error decrypting data: %s", err)
	}

	var sd serviceData
	err = json.Unmarshal(dec, &sd)

	if err != nil {
		fmt.Printf("data is decrypted successfully but the content is corrupted: %s\ninitializing empty store...", err)
		m.data = make(serviceData)
	} else {
		m.data = sd
	}

	return nil
}

func (m *Manager) save() error {
	data, err := json.Marshal(m.data)
	if err != nil {
		return fmt.Errorf("error marshaling data: %s", err)
	}

	enc, err := m.enc.Encrypt(data)
	if err != nil {
		return fmt.Errorf("error encrypting data: %s", err)
	}

	err = m.store.Save(enc)
	if err != nil {
		return fmt.Errorf("error saving data: %s", err)
	}

	return nil
}

func (m *Manager) cmd(line string) {
	var args []string
	var argsLine string

	line = strings.Trim(line, " \n\t")

	cmdRunes, rest := wsSplit([]rune(line))
	cmd := string(cmdRunes)

	if cmd == "" {
		return
	}

	if rest == nil {
		args = make([]string, 0)
		argsLine = ""
	} else {
		argsLine = string(rest)
		args = exprWhiteSpace.Split(argsLine, -1)
	}

	if handler, ok := m.handlers[cmd]; ok {
		handler(cmd, argsLine, args...)
	} else {
		term.Errorf("Unknown command: %s\n", cmd)
	}
}

// Run starts the manager loop
func (m *Manager) Run() {
	fmt.Println("Loading remote data...")
	err := m.load()
	if err != nil {
		fmt.Println(err)
		return
	}
	m.setPrompt()
	m.cmdLoop()
}
