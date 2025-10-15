# QOwnNotes command-line snippet manager changelog

## v0.6.2

- Refactor the code to fix deprecations and issues and to make it cleaner and easier to maintain
- Update dependencies

## v0.6.1

- The `--atuin` flag now also add commands to [Atuin](https://atuin.sh/) for multi-line commands
  (for [#15](https://github.com/qownnotes/qc/issues/15))
- The `--color` flag now shows the command description in a calmer green, instead of red
  (for [#16](https://github.com/qownnotes/qc/issues/16))

## v0.6.0

- Add support for storing commands in [Atuin](https://atuin.sh/) on execution
  when using the `--atuin` flag (for [#15](https://github.com/qownnotes/qc/issues/15))
  - This only works for single-line commands
- Update dependencies

## v0.5.1

- The last selected command is now only stored when there actually was a command selected and
  the dialog wasn't quit without selecting a command (for [#9](https://github.com/qownnotes/qc/issues/9))

## v0.5.0

- The last executed command is now stored and can be executed again via `qc exec --last`
  (for [#9](https://github.com/qownnotes/qc/issues/9))
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
