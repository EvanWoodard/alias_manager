package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/EvanWoodard/alias_manager/server"
)

const (
	zshrc    = ".zshrc"
	zshAlias = "zsh_alias"
)

var (
	homeDir, _ = os.UserHomeDir()
	aliasPath  = filepath.Join(homeDir, ".zsh", "am")
)

func main() {
	fmt.Println("Hi! I'm your alias manager. You can call me AL! Give me a moment to set things up!")

	setupRC()

	ctx, cancel := context.WithCancel(context.Background())

	srv := server.New(aliasPath, zshAlias)

	args := os.Args[1:]
	if len(args) > 0 {
		srv.RunCmd(ctx, args[0])
		return
	}
	srv.Run(ctx)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

		// Block until we receive our signal.
	<-c
	cancel()

	os.Exit(0)
}

func setupRC() {
	z := filepath.Join(homeDir, zshrc)
	file, err := os.OpenFile(z, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(z)
			if err != nil {
				fmt.Println("Uh-Oh, could not create a .zshrc file...")
				os.Exit(1)
			}
		}
	}

	contents, err := os.ReadFile(z)
	if err != nil {
		fmt.Println(err)
	}
	if !strings.Contains(string(contents), server.ImportAlias) {
		fmt.Println(".zshrc file does not contain alias import, adding now...")
		_, err = file.WriteString(server.ImportAlias)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(".zshrc file already contains alias import, skipping...")
	}
}
