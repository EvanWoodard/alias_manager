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

type aliasServer struct{}

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

	file, err := os.OpenFile(filepath.Join(aliasPath, zshAlias), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			a.log("Alias file is missing, please restart the alias server to re-initialize it")
			a.log("Shutting down...")
			os.Exit(1)
		}
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("alias %s=\"%s\"\n", name, fn))
	if err != nil {
		a.log(err.Error())
		return
	}

	cmd := exec.Command("bash", "-c", "source ~/.zsh/am/zsh_alias")
	a.log("Re-sourcing your aliases, you should not have to restart your terminal to see it applied.")
	err = cmd.Run()
	if err != nil {
		a.log("Something went wrong while sourcing aliases, your alias was saved, but not applied. Restart your terminal to use your new alias")
		a.log(err.Error())
	}
}

func (a *aliasServer) listAliases() {
	a.log("Listing Aliases")

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

			a.writeToUser(fmt.Sprintf("%s : %s", cmdList[0], cmdList[1]))
		}
	}

	if err != io.EOF {
		a.log(fmt.Sprintf("Failed!: %v\n", err))
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
