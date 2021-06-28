package main

const (
	importAlias    = "\n# This imports the aliases from alias manager\nif [ -f ~/.zsh/am/zsh_alias ]; then\n\tsource ~/.zsh/am/zsh_alias\nelse\n\tprint \"404: ~/.zsh/am/zsh_alias not found.\"\nfi"
	defaultAliases = "\n# These aliases are default and can't be removed. If you want them removed, have a word with Evan\nalias ga=\"git add --all\"\nalias gc=\"git commit -m\"\nalias gp=\"git push\"\n\n# Here begins your custom aliases\n"
)
