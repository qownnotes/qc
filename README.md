# QOwnNotes command-line snippet manager

[GitHub](https://github.com/qownnotes/qc) |
[Changelog](https://github.com/qownnotes/qc/blob/main/CHANGELOG.md)

You can use the QOwnNotes command-line snippet manager to execute command snippets stored
in notes in [QOwnNotes](https://www.qownnotes.org/) from the command line.

You can use notes with a special tag to store command snippets, which you can
execute from the command-line snippet manager.

For more information on how to add commands please see
[Command-line Snippet Manager](https://www.qownnotes.org/getting-started/command-line-snippet-manager.html).

The QOwnNotes command-line snippet manager is based on the
[pet CLI Snippet Manager](https://github.com/knqyf263/pet).

## Dependencies

- [fzf](https://github.com/junegunn/fzf) or [peco](https://github.com/peco/peco)

## Usage

```
Usage:
qc [command]

Available Commands:
completion  generate the autocompletion script for the specified shell
configure   Edit config file
exec        Run the selected commands
help        Help about any command
search      Search snippets
version     Print the version number

Flags:
--config string   config file (default is $HOME/.config/qc/config.toml)
--debug           debug mode
-h, --help        help for qc

Use "qc [command] --help" for more information about a command.
```
