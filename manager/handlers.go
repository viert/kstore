package manager

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/viert/kstore/term"
)

func (m *Manager) setupHandlers() {
	m.handlers = make(map[string]cmdHandler)
	m.handlers["exit"] = m.doExit
	m.handlers["save"] = m.doSave
	m.handlers["import"] = m.doImport
	m.handlers["export"] = m.doExport
	m.handlers["list"] = m.doList
	m.handlers["ls"] = m.doList
	m.handlers["get"] = m.doGet
	m.handlers["getpass"] = m.doGet
	m.handlers["set"] = m.doSet
	m.handlers["setpass"] = m.doSet
	m.handlers["delete"] = m.doDelete
	m.handlers["remove"] = m.doDelete
	m.handlers["del"] = m.doDelete
	m.handlers["rm"] = m.doDelete
}

func (m *Manager) doExit(name string, argsLine string, args ...string) {
	m.stopped = true
}

func (m *Manager) doSave(name string, argsLine string, args ...string) {
	err := m.save()
	if err != nil {
		term.Errorf("%s\n", err)
	}
}

func (m *Manager) doImport(name string, argsLine string, args ...string) {
	if len(args) < 1 {
		term.Errorf("import command requires a file name to import from\n")
		return
	}

	filename := args[0]

	f, err := os.Open(filename)
	if err != nil {
		term.Errorf("error opening file %s: %s\n", filename, err)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		term.Errorf("error reading file %s: %s\n", filename, err)
		return
	}

	sd := make(serviceData)
	err = json.Unmarshal(data, &sd)
	if err != nil {
		term.Errorf("error parsing file %s: %s\n", filename, err)
		return
	}

	added := 0
	skipped := 0
	for k, v := range sd {
		if _, found := m.data[k]; found {
			term.Warnf("Service %s already exists, skipping\n", k)
			skipped++
			continue
		}
		m.data[k] = v
		added++
	}

	term.Successf("%d items imported, %d items skipped\n", added, skipped)
	term.Warnf("The imported data is not persistent yet, don't forget to **save** it.\n")
}

func (m *Manager) doExport(name string, argsLine string, args ...string) {
	if len(args) < 1 {
		term.Errorf("Use export <filename> to export data to a json file\n")
	}

	data, err := m.data.dump()
	if err != nil {
		term.Errorf("Error marshaling data: %s\n", err)
		return
	}

	filename := args[0]
	err = ioutil.WriteFile(filename, data, os.FileMode(0644))
	if err != nil {
		term.Errorf("Error writing file %s: %s\n", filename, err)
		return
	}

	term.Successf("Data has been successfully exported to file %s\n", filename)
}

func (m *Manager) doList(name string, argsLine string, args ...string) {
	if len(m.data) > 0 {
		names := make([]string, len(m.data))
		i := 0
		for k := range m.data {
			names[i] = k
			i++
		}
		sort.Strings(names)

		for _, name := range names {
			fmt.Println(name)
		}
	} else {
		term.Errorf("Empty list\n")
	}
}

func (m *Manager) doGet(name string, argsLine string, args ...string) {
	if len(args) < 1 {
		term.Errorf("%s command requires a service name\n", name)
		return
	}

	serviceName := args[0]
	if item, found := m.data[serviceName]; found {
		switch name {
		case "getpass":
			fmt.Println(item.Password)
		default:
			fmt.Printf("Service: %s\n", item.Name)
			if item.Username != "" {
				fmt.Printf("Username: %s\n", item.Username)
			}
			if item.Password != "" {
				fmt.Printf("Password: %s\n", item.Password)
			}
			if item.Comment != "" {
				fmt.Printf("Comment: %s\n", item.Comment)
			}
			if item.URL != "" {
				fmt.Printf("URL: %s\n", item.URL)
			}
			fmt.Println()
		}
	} else {
		term.Errorf("Service %s not found\n", serviceName)
	}
}

func (m *Manager) doSet(name string, argsLine string, args ...string) {
	if len(args) < 1 {
		term.Errorf("%s command requires a service name\n", name)
		return
	}

	serviceName := args[0]

	switch name {
	case "setpass":
		si, found := m.data[serviceName]
		if !found {
			term.Errorf(
				"Service %s not found. If you want to create it, use \"set\" command instead of setpass\n",
				serviceName,
			)
			return
		}

		pwd, err := getString("Password: ")
		if err != nil {
			return
		}
		si.Password = pwd
		term.Successf("Password updated. Don't forget to **save** the result.\n")
	default:
		si := &ServiceInfo{Name: serviceName}
		si.Username, _ = getString("Username: ")
		si.Password, _ = getString("Password: ")
		si.Comment, _ = getString("Comment: ")
		si.URL, _ = getString("URL: ")
		m.data[serviceName] = si
		term.Successf("Service %s created. Don't forget to **save** the result.\n", serviceName)
	}
}

func (m *Manager) doDelete(name string, argsLine string, args ...string) {
	if len(args) < 1 {
		term.Errorf("%s command requires a service name\n", name)
		return
	}

	serviceName := args[0]
	_, found := m.data[serviceName]
	if !found {
		term.Errorf("Service %s not found.", serviceName)
		return
	}

	delete(m.data, serviceName)
	term.Successf("Service %s removed. Don't forget to **save** the result.\n", serviceName)
}

func getString(prompt string) (string, error) {
	fmt.Print(prompt)
	rd := bufio.NewReader(os.Stdin)
	data, err := rd.ReadString('\n')
	data = strings.TrimSpace(data)
	return data, err
}
