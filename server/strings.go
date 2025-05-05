package server

const (
	importAlias    = "\n# This imports the aliases from alias manager\nif [ -f ~/.zsh/am/zsh_alias ]; then\n\tsource ~/.zsh/am/zsh_alias\nelse\n\tprint \"404: ~/.zsh/am/zsh_alias not found.\"\nfi\ntrap 'source ~/.zsh/am/zsh_alias' USR1\n"
	defaultAliases = "# These are some default aliases, feel free to remove them\nalias ga=git add --all\nalias gc=git commit -m\nalias gp=git push\n\n# Here begins your custom aliases\n"
)
