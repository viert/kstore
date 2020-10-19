package manager

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/chzyer/readline"
)

type completeFunc func([]rune) ([][]rune, int)

type cliCompleter struct {
	commands   []string
	completers map[string]completeFunc
}

var (
	exprWhiteSpace = regexp.MustCompile(`\s+`)
)

func wsSplit(line []rune) ([]rune, []rune) {
	sline := string(line)
	tokens := exprWhiteSpace.Split(sline, 2)
	if len(tokens) < 2 {
		return []rune(tokens[0]), nil
	}
	return []rune(tokens[0]), []rune(tokens[1])
}

func toRunes(src []string) [][]rune {
	dst := make([][]rune, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = []rune(src[i])
	}
	return dst
}

func staticCompleter(variants []string) completeFunc {
	sort.Strings(variants)
	return func(line []rune) (newLine [][]rune, length int) {
		ll := len(line)
		sr := make([]string, 0)
		for _, variant := range variants {
			if strings.HasPrefix(variant, string(line)) {
				sr = append(sr, variant[ll:])
			}
		}
		return toRunes(sr), ll
	}
}

func newCliCompleter(commands []string) *cliCompleter {
	c := &cliCompleter{commands, make(map[string]completeFunc)}
	return c
}

func (c *cliCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	postfix := line[pos:]
	result, length := c.complete(line[:pos])
	if len(postfix) > 0 {
		for i := 0; i < len(result); i++ {
			result[i] = append(result[i], postfix...)
		}
	}
	return result, length
}

func (c *cliCompleter) complete(line []rune) (newLine [][]rune, length int) {
	cmd, args := wsSplit(line)
	if args == nil {
		return c.completeCommand(cmd)
	}

	if handler, found := c.completers[string(cmd)]; found {
		return handler(args)
	}

	return [][]rune{}, 0
}

func (c *cliCompleter) completeCommand(line []rune) (newLine [][]rune, length int) {
	sr := make([]string, 0)
	for _, cmd := range c.commands {
		if strings.HasPrefix(cmd, string(line)) {
			sr = append(sr, cmd[len(line):]+" ")
		}
	}
	sort.Strings(sr)
	return toRunes(sr), len(line)
}

func (m *Manager) nameCompleter() completeFunc {
	return func(line []rune) (newLine [][]rune, length int) {
		sr := make([]string, 0)
		for name := range m.data {
			if strings.HasPrefix(name, string(line)) {
				sr = append(sr, name[len(line):])
			}
		}
		sort.Strings(sr)
		return toRunes(sr), len(line)
	}
}

func (m *Manager) setupReadline() error {

	commands := make([]string, len(m.handlers))

	i := 0
	for cmd := range m.handlers {
		commands[i] = cmd
		i++
	}

	cc := newCliCompleter(commands)
	nc := m.nameCompleter()
	cc.completers["get"] = nc
	cc.completers["getpass"] = nc
	cc.completers["set"] = nc
	cc.completers["setpass"] = nc
	cc.completers["delete"] = nc
	cc.completers["remove"] = nc
	cc.completers["rm"] = nc
	cc.completers["del"] = nc

	readlineConfig := &readline.Config{
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		AutoComplete:      cc,
	}

	rl, err := readline.NewEx(readlineConfig)
	if err != nil {
		fmt.Printf("Error creating readline instance: %s\n", err)
		return err
	}
	m.rl = rl

	return nil
}
