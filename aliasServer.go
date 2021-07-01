package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type cliServer interface {
	Shutdown(context.Context)
}

type aliasServer struct {
	aliases map[string]string
}

func startAliasServer() cliServer {
	aliasServer := aliasServer{}

	aliasServer.Start()

	return &aliasServer
}

func (a *aliasServer) Shutdown(ctx context.Context) {
	a.writeToUser("")
	a.log("Shutting down...")
}

func (a *aliasServer) Start() {
	a.aliases = make(map[string]string)
	go a.LoopInput()
}

func (a *aliasServer) LoopInput() {
	r := bufio.NewReader(os.Stdin)
	a.writeToUser("")
	a.writeToUser("Alias Shell")
	a.writeToUser("---------------------")

	for {
		a.promptUser("->")
		cmd := a.readLine(r)

		switch cmd {
		case "new", "create", "n":
			a.promptUser("New alias name:")
			name := a.readLine(r)
			a.promptUser("Alias function:")
			fn := a.readLine(r)

			a.createAlias(name, fn)
		case "list", "l":
			a.listAliases()
		case "remove", "r":
			a.promptUser("Alias to remove:")
			alias := a.readLine(r)
			a.removeAlias(alias)
		case "hi", "hello":
			a.log("hello, Yourself")
		}
	}
}

func (a *aliasServer) readLine(r *bufio.Reader) string {
	text, _ := r.ReadString('\n')
	// convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\"", "\\\"", -1)

	return text
}

func (a *aliasServer) createAlias(name, fn string) {
	a.log(fmt.Sprintf("Creating your alias, sit tight! Alias: %s:%s", name, fn))

	a.aliases[name] = fn
	a.writeAliases()
}

func (a *aliasServer) listAliases() {
	a.log("Listing Aliases")
	a.checkAliases()
	for cmd, fn := range a.aliases {
		a.writeToUser(fmt.Sprintf("%s : %s", cmd, fn))
	}
}

func (a *aliasServer) removeAlias(alias string) {
	delete(a.aliases, alias)
	a.writeAliases()
}

func (a *aliasServer) checkAliases() {
	file, err := os.Open(filepath.Join(aliasPath, zshAlias))
	if err != nil {
		a.log("Could not find the alias list file. Restart the alias server to re-initialize it")
		return
	}
	defer file.Close()

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)

	var line string
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		}

		if len(line) < 5 {
			continue
		}

		if line[:5] == "alias" {
			cmd := line[6:]
			cmd = strings.Replace(cmd, "\n", "", -1)
			cmd = strings.Replace(cmd, "\"", "", -1)
			cmdList := strings.Split(cmd, "=")

			a.aliases[cmdList[0]] = cmdList[1]
		}
	}

	if err != io.EOF {
		a.log(fmt.Sprintf("Failed!: %v\n", err))
	}
}

func (a *aliasServer) writeAliases() {
	aliasStr := ""
	for cmd, fn := range a.aliases {
		aliasStr += fmt.Sprintf("alias %s=\"%s\"\n", cmd, fn)
	}

	file, err := os.OpenFile(filepath.Join(aliasPath, zshAlias), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			a.log("Could not find the alias list file. Restart the alias server to re-initialize it")
		}
		return
	}

	a.log("Writing aliases to file")
	_, err = file.WriteString(aliasStr)
	if err != nil {
		a.log(err.Error())
	}

	a.log("Getting terminal instances")
	// grep := exec.Command("ps", "ax", "|", "grep", "/bin/zsh")
	grep := exec.Command("grep", "zsh")
	ps := exec.Command("ps", "ax")

	psPipe, _ := ps.StdoutPipe()
	defer psPipe.Close()
	grep.Stdin = psPipe
	ps.Start()

	res, _ := grep.Output()

	lines := strings.Split(strings.TrimSpace(string(res)), "\n")
	var pids []string
	for _, line := range lines {
		lineArr := strings.Split(line, " ")
		pids = append(pids, lineArr[0])
	}

	for _, pid := range pids {
		usrCmd := exec.Command("kill", "-USR1", pid)
		usrCmd.Start()
	}
}

func (a *aliasServer) log(msg string) {
	fmt.Printf("AL: %s\n", msg)
}

func (a *aliasServer) writeToUser(msg string) {
	fmt.Printf("%s\n", msg)
}

func (a *aliasServer) promptUser(msg string) {
	fmt.Printf("%s ", msg)
}
