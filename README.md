# QOwnNotes command-line snippet manager

[GitHub](https://github.com/qownnotes/qc) |
[Changelog](https://github.com/qownnotes/qc/blob/main/CHANGELOG.md) |
[Releases](https://github.com/qownnotes/qc/releases)

You can use the **QOwnNotes command-line snippet manager** to **execute command snippets** stored
in **notes** in [QOwnNotes](https://www.qownnotes.org/) from the command line.

![qc](qc.png)

You can use **notes with a special tag** (`commands` by default) to store **command snippets**, which you can
**execute from the command-line snippet manager**.

![commands](commands.png)

For more information on **how to add commands and configuration** please see
[Command-line Snippet Manager](https://www.qownnotes.org/getting-started/command-line-snippet-manager.html).

The QOwnNotes command-line snippet manager is based on the wonderful
[pet CLI Snippet Manager](https://github.com/knqyf263/pet).

## Installation

Visit the [latest release page](https://github.com/qownnotes/qc/releases/latest)
and download the version you need.

If you have [jq](https://stedolan.github.io/jq) installed you can also use this snippet
to download and install for example the latest Linux AMD64 AppImage to `/usr/local/bin/qc`:

```bash
curl https://api.github.com/repos/qownnotes/qc/releases/latest | \
jq '.assets[] | select(.browser_download_url | endswith("_linux_amd64.tar.gz")) | .browser_download_url' | \
xargs curl -Lo /tmp/qc.tar.gz && \
tar xfz /tmp/qc.tar.gz -C /tmp && \
rm /tmp/qc.tar.gz && \
sudo mv /tmp/qc /usr/local/bin/qc && \
/usr/local/bin/qc version
```

### macOS / Homebrew

You can use homebrew on macOS to install qc.

```bash
brew install qownnotes/qc/qc
```

If you receive an error (`Error: qownnotes/qc/qc 64 already installed`) during `brew upgrade`,
try the following command:

```bash
brew unlink qc && brew uninstall qc
rm -rf /usr/local/Cellar/qc/64
brew install qownnotes/qc/qc
```

## Dependencies

[fzf](https://github.com/junegunn/fzf) (fuzzy search) or [peco](https://github.com/peco/peco)
(older, but more likely to be installed by default) need to be installed to search
for commands on the command-line.

By default `fzf` is used for searching, but you can use `peco` by setting it with `qc configure`.

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

## Configuration

Run `qc configure`.

```toml
[General]
  editor = "vim"            # your favorite text editor
  column = 40               # column size for list command
  selectcmd = "fzf"         # selector command for edit command (fzf or peco)
  sortby = ""               # specify how snippets get sorted (recency (default), -recency, description, -description, command, -command, output, -output)

[QOwnNotes]
  token = "SECRET"          # your QOwnNotes API token
  websocket_port = 22222    # websocket port in QOwnNotes
```

## Shell completion

You can generate shell completion code for your shell with `qc completion <shell>`.

For example for the Fish shell you can use:

```bash
qc completion fish > ~/.config/fish/completions/qc.fish
```
