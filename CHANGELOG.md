# QOwnNotes command-line snippet manager changelog

## v0.5.0
- Neovim is now also used to edit the config file 

## v0.4.0
- Add support for note folder switching (for [#5](https://github.com/qownnotes/qc/issues/5))
  - Use `qc switch` to get a list of note folders to select which note folder to switch to
  - Use `qc switch -f <note-folder-id>` to make QOwnNotes switch to another note folder instantly
  - Needs OwnNotes at least at version 22.7.1

## v0.3.2
- Add Homebrew tap for qc (`brew install qownnotes/qc/qc`)

## v0.3.0
- Enable sorting of snippets via settings and allow sorting case-insensitively

## v0.2.0
- Cache snippets in case QOwnNotes is not running
- Don't throw an error if selectCmd was exited with an error code (e.g. by `Ctrl + C`)

## v0.1.0
- Support for fetching snippets from QOwnNotes via websocket
- Searching in snippets
- Executing snippets
- Configuring the application
- Generating autocompletion scripts for shells
