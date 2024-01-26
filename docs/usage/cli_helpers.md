---
title: Shell completion
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 1000
description: shell autocompletion for the go-swagger CLI
---
# CLI helpers

## Bash Completion

Bash completion is supported and can be activated as follows:

```bash
source ./cmd/swagger/completion/swagger.bash-completion
```

Note that this does require you already setup bash completion,
which can be done in 2 simple steps:

1) install `bash-completion` using your favorite package manager;
2) run `source /etc/bash_completion` in bash;

## Zsh Completion

Zsh completion is supported and can be copied/soft-linked from:

```zsh
./cmd/swagger/completion/swagger.zsh-completion
```

In case you're new to adding auto-completion to zsh completion,
here is how you could enable swagger's zsh completion step by step:

1) create a folder used to store your completions (eg. `$HOME/.zsh/completion`);
2) append the following to your `$HOME/.zshrc` file:

```zsh
# add auto-completion directory to zsh's fpath
fpath=($HOME/.zsh/completion $fpath)

# compsys initialization
autoload -U compinit
compinit
```

3) copy/soft-link `./cmd/swagger/completion/swagger.zsh-completion` to `$HOME/.zsh/completion/_swagger`;
