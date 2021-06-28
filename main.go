package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
	fmt.Println("Hi! I'm your alias manager, give me a moment to set things up!")

	setupRC()
	setupAliasFile()
	// startupAliasServer()
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

	contents, err := ioutil.ReadFile(z)
	if err != nil {
		fmt.Println(err)
	}
	if !strings.Contains(string(contents), importAlias) {
		fmt.Println(".zshrc file does not contain alias import, adding now...")
		_, err = file.WriteString(importAlias)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(".zshrc file already contains alias import, skipping...")
	}
}

func setupAliasFile() {
	err := os.MkdirAll(aliasPath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.OpenFile(filepath.Join(aliasPath, zshAlias), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(filepath.Join(aliasPath, zshAlias))
			if err != nil {
				fmt.Println("Uh-Oh, could not create a zsh_alias file...")
				os.Exit(1)
			}
		}
	}

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
