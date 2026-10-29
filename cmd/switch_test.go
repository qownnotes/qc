package cmd

import (
	"errors"
	"os/exec"
	"runtime"
	"testing"

	"github.com/qownnotes/qc/config"
	"github.com/qownnotes/qc/entity"
	"gopkg.in/alessio/shellescape.v1"
)

func TestSelectNoteFolder(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test selectors use POSIX shell commands")
	}
	original := config.Conf
	t.Cleanup(func() { config.Conf = original })
	folders := []entity.NoteFolderInfo{
		{Id: 7, Name: "Duplicate"},
		{Id: 1, Name: "Notes"},
		{Id: 8, Name: "Duplicate"},
		{Id: 0, Name: "Tabs\tand\nlines"},
	}
	for _, tt := range []struct {
		name      string
		command   string
		id        int
		selected  bool
		wantError bool
		interrupt bool
	}{
		{"current first", "head -n 1", 1, true, false, false},
		{"duplicate name", "sed -n '3p'", 8, true, false, false},
		{"zero id and sanitized name", "tail -n 1", 0, true, false, false},
		{"empty selection", "cat >/dev/null", 0, false, false, false},
		{"no match", "exit 1", 0, false, false, false},
		{"Ctrl+C status", "exit 130", 0, false, true, true},
		{"SIGINT", "kill -INT $$", 0, false, true, true},
		{"selector failure", "exit 2", 0, false, true, false},
		{"unknown output", "printf 'unknown\\n'", 0, false, true, false},
		{"multiple selections", "head -n 2", 0, false, true, false},
		{"missing configuration", "", 0, false, true, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			config.Conf.General.SelectCmd = tt.command
			id, selected, err := selectNoteFolder(folders, 1)
			if id != tt.id || selected != tt.selected || (err != nil) != tt.wantError {
				t.Fatalf("got (%d, %v, %v), expected (%d, %v), error=%v",
					id, selected, err, tt.id, tt.selected, tt.wantError)
			}
			if errors.Is(err, ErrSigInt) != tt.interrupt {
				t.Fatalf("expected interrupted=%v, got %v", tt.interrupt, err)
			}
		})
	}
}

func TestSelectNoteFolderEmpty(t *testing.T) {
	if _, selected, err := selectNoteFolder(nil, 1); err == nil || selected {
		t.Fatalf("expected no folders error, got selected=%v, err=%v", selected, err)
	}
}

func TestSelectNoteFolderFzf(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test command uses POSIX quoting")
	}
	fzf, err := exec.LookPath("fzf")
	if err != nil {
		t.Skip("fzf is not installed")
	}
	original := config.Conf
	t.Cleanup(func() { config.Conf = original })
	config.Conf.General.SelectCmd = shellescape.Quote(fzf) + " --filter=Notes"
	id, selected, err := selectNoteFolder([]entity.NoteFolderInfo{
		{Id: 7, Name: "TU"},
		{Id: 1, Name: "Notes"},
	}, 1)
	if err != nil || !selected || id != 1 {
		t.Fatalf(
			"expected fzf to return the original row with id 1, got (%d, %v, %v)",
			id,
			selected,
			err,
		)
	}
}

func TestSwitchNoteFolderInvalidID(t *testing.T) {
	original := config.Flag
	t.Cleanup(func() { config.Flag = original })
	for _, id := range []string{"not-a-number", "-1"} {
		config.Flag.Query = id
		if err := switchNoteFolder(nil, nil); err == nil {
			t.Fatalf("expected invalid id %q to be rejected", id)
		}
	}
}
