package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

// Conf is global config variable
var Conf Config

// Config is a struct of config
type Config struct {
	General   GeneralConfig
	QOwnNotes QOwnNotesConfig
}

// QOwnNotesConfig is a struct of config for QOwnNotes
type QOwnNotesConfig struct {
	Token         string `toml:"token"`
	WebSocketPort int    `toml:"websocket_port"`
}

// Flag is global flag variable
var Flag FlagConfig

// FlagConfig is a struct of flag
type FlagConfig struct {
	Debug     bool
	Query     string
	Command   bool
	Atuin     bool
	FilterTag string
	Color     bool
	Delimiter string
	Last      bool
}

// GeneralConfig is a struct of general config
type GeneralConfig struct {
	Editor    string `toml:"editor"`
	Column    int    `toml:"column"`
	SelectCmd string `toml:"selectcmd"`
	SortBy    string `toml:"sortby"`
}

// Load loads a config toml
func (cfg *Config) Load(file string) error {
	_, err := os.Stat(file)
	if err == nil {
		_, err := toml.DecodeFile(file, cfg)
		if err != nil {
			return err
		}
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	cfg.General.Editor = os.Getenv("EDITOR")
	if cfg.General.Editor == "" && runtime.GOOS != "windows" {
		if isCommandAvailable("sensible-editor") {
			cfg.General.Editor = "sensible-editor"
		} else if isCommandAvailable("nvim") {
			cfg.General.Editor = "nvim"
		} else {
			cfg.General.Editor = "vim"
		}
	}
	cfg.General.Column = 40
	cfg.General.SelectCmd = "fzf"
	cfg.QOwnNotes.WebSocketPort = 22222

	return toml.NewEncoder(f).Encode(cfg)
}

// GetDefaultConfigDir returns the default config directory
func GetDefaultConfigDir() (dir string, err error) {
	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", "qc")
		}
		dir = filepath.Join(dir, "qc")
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config", "qc")
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("cannot create directory: %v", err)
	}
	return dir, nil
}

func expandPath(s string) string {
	if len(s) >= 2 && s[0] == '~' && os.IsPathSeparator(s[1]) {
		if runtime.GOOS == "windows" {
			s = filepath.Join(os.Getenv("USERPROFILE"), s[2:])
		} else {
			s = filepath.Join(os.Getenv("HOME"), s[2:])
		}
	}
	return os.Expand(s, os.Getenv)
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
