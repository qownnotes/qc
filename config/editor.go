package config

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ErrNoEditor is returned when no usable text editor could be found on the system.
var ErrNoEditor = errors.New(
	"no usable text editor was found: set the VISUAL or EDITOR environment variable, " +
		"or the general.editor option in the qc config file")

// editorEnvVars are the environment variables that are consulted, in order, to
// determine the editor preferred by the user.
var editorEnvVars = []string{"VISUAL", "EDITOR"}

// editorCandidates returns the editor commands that are probed, in order of
// preference, when neither the config file nor the environment provide a
// usable editor.
func editorCandidates() []string {
	if runtime.GOOS == "windows" {
		return []string{"nvim", "vim", "notepad++", "notepad"}
	}
	return []string{
		"sensible-editor", // Debian/Ubuntu wrapper honouring the user preference
		"editor",          // Debian alternatives symlink
		"nvim", "vim", "vi",
		"nano", "micro", "hx", "helix", "emacs", "mcedit",
	}
}

// IsEditorAvailable reports whether the given editor command line can actually
// be executed on this system. The command may contain arguments, for example
// "code --wait".
func IsEditorAvailable(command string) bool {
	command = strings.TrimSpace(command)
	if command == "" {
		return false
	}
	// Handles an unquoted path that contains spaces but no arguments, such as
	// `C:\Program Files\Editor\editor.exe`.
	if _, err := exec.LookPath(command); err == nil {
		return true
	}
	executable := editorExecutable(command)
	if executable == "" || executable == command {
		return false
	}
	_, err := exec.LookPath(executable)
	return err == nil
}

// editorExecutable extracts the executable part of an editor command line.
func editorExecutable(command string) string {
	command = strings.TrimSpace(command)
	if command == "" {
		return ""
	}
	if quote := command[0]; quote == '"' || quote == '\'' {
		if end := strings.IndexByte(command[1:], quote); end >= 0 {
			return command[1 : end+1]
		}
	}
	return strings.Fields(command)[0]
}

// DetectEditor determines the editor to use by looking at the VISUAL and
// EDITOR environment variables first and by probing a list of well known
// editors afterwards. It returns an empty string when nothing was found.
func DetectEditor() string {
	for _, name := range editorEnvVars {
		if editor := strings.TrimSpace(os.Getenv(name)); IsEditorAvailable(editor) {
			return editor
		}
	}
	for _, candidate := range editorCandidates() {
		if IsEditorAvailable(candidate) {
			return candidate
		}
	}
	return ""
}

// ResolveEditor returns a usable editor command. The configured editor is used
// when it is installed, otherwise an editor is detected on the fly so that a
// stale or empty config entry does not make qc fail. The returned bool reports
// whether the configured editor had to be replaced by a detected one.
func ResolveEditor(configured string) (editor string, replaced bool, err error) {
	if IsEditorAvailable(configured) {
		return strings.TrimSpace(configured), false, nil
	}
	detected := DetectEditor()
	if detected == "" {
		return "", false, ErrNoEditor
	}
	return detected, strings.TrimSpace(configured) != "", nil
}
