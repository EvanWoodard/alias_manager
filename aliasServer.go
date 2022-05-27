package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	a.log("Shutting down...")
	os.Exit(0)
}

func (a *aliasServer) Start() {
	a.aliases = make(map[string]string)
	a.setupAliasFile()
	a.checkAliases()
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
		case "exit", "close", "q":
			a.log("Bye")
			a.Shutdown(nil)
		case "hi", "hello":
			a.log("hello, Yourself")
		case "help", "h":
			a.writeToUser("Commands:")
			a.writeToUser("new (create, n): Creates new alias")
			a.writeToUser("list (l): Lists created aliases")
			a.writeToUser("remove (r): Removes alias")
			a.writeToUser("exit (close, q): exit Al")
		}
	}
}

func (a *aliasServer) readLine(r *bufio.Reader) string {
	text, _ := r.ReadString('\n')

	text = strings.Replace(text, "\n", "", -1)
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
			cmd = strings.Replace(cmd, "\\\"", "\"", -1)
			cmd = strings.Replace(cmd, "'", "", -1)
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
		aliasStr += fmt.Sprintf("alias %s='%s'\n", cmd, fn)
	}

	file, err := os.OpenFile(filepath.Join(aliasPath, zshAlias), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			a.log("Could not find the alias list file. Restart the alias server to re-initialize it")
		}
		return
	}

	// This will clear the contents of the file by setting the filesize to 0
	if err := file.Truncate(0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	a.log("Writing aliases to file")
	_, err = file.WriteString(aliasStr)
	if err != nil {
		a.log(err.Error())
	}
	
	a.updateAllRunningShells()
}

func (a *aliasServer) updateAllRunningShells() {
	a.log("Getting terminal instances")
	grep := exec.Command("grep", "/bin/zsh \\| -zsh")
	ps := exec.Command("ps", "ax")

	psPipe, _ := ps.StdoutPipe()

	defer psPipe.Close()
	grep.Stdin = psPipe
	ps.Start()

	res, _ := grep.Output()

	lines := strings.Split(strings.TrimSpace(string(res)), "\n")
	var pids []string
	for _, line := range lines {
		if strings.Contains(line, "grep") {
			continue
		}
		line = strings.TrimSpace(line)
		lineArr := strings.Split(line, " ")
		pids = append(pids, lineArr[0])
	}

	for _, pid := range pids {
		usrCmd := exec.Command("kill", "-USR1", pid)
		usrCmd.Start()
	}
}

func (a *aliasServer) setupAliasFile() {
	err := os.MkdirAll(aliasPath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = os.OpenFile(filepath.Join(aliasPath, zshAlias), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(filepath.Join(aliasPath, zshAlias))
			if err != nil {
				fmt.Println("Uh-Oh, could not create a zsh_alias file...")
				os.Exit(1)
			}

			a.writeDefaultAliases(file)
		}
	}
}

func (a *aliasServer) writeDefaultAliases(file *os.File) {
	contents, err := ioutil.ReadFile(filepath.Join(aliasPath, zshAlias))
	if err != nil {
		fmt.Println(err)
	}
	if !strings.Contains(string(contents), defaultAliases) {
		fmt.Println("Alias file does not contain default aliases, adding now...")
		_, err = file.WriteString(defaultAliases)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Alias file already contains default aliases, skipping...")
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
