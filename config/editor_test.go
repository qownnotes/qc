package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// fakeEditor creates an executable named name in dir.
func fakeEditor(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("cannot create fake editor: %v", err)
	}
	return path
}

func TestEditorExecutable(t *testing.T) {
	tests := []struct {
		command string
		want    string
	}{
		{"", ""},
		{"vim", "vim"},
		{"  vim  ", "vim"},
		{"code --wait", "code"},
		{`"/opt/my editor/ed" --wait`, "/opt/my editor/ed"},
		{`'/opt/my editor/ed'`, "/opt/my editor/ed"},
	}
	for _, test := range tests {
		if got := editorExecutable(test.command); got != test.want {
			t.Errorf("editorExecutable(%q) = %q, want %q", test.command, got, test.want)
		}
	}
}

func TestIsEditorAvailable(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell based fake editors are not portable to Windows")
	}
	dir := t.TempDir()
	fakeEditor(t, dir, "qc-test-editor")
	t.Setenv("PATH", dir)

	if !IsEditorAvailable("qc-test-editor") {
		t.Error("expected the fake editor to be reported as available")
	}
	if !IsEditorAvailable("qc-test-editor --wait -n") {
		t.Error("expected an editor with arguments to be reported as available")
	}
	if IsEditorAvailable("qc-test-missing-editor") {
		t.Error("expected a missing editor to be reported as unavailable")
	}
	if IsEditorAvailable("   ") {
		t.Error("expected a blank editor to be reported as unavailable")
	}
}

func TestDetectEditorPrefersEnvironment(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell based fake editors are not portable to Windows")
	}
	dir := t.TempDir()
	fakeEditor(t, dir, "qc-visual-editor")
	fakeEditor(t, dir, "qc-env-editor")
	fakeEditor(t, dir, "vim")
	t.Setenv("PATH", dir)

	t.Setenv("VISUAL", "qc-visual-editor")
	t.Setenv("EDITOR", "qc-env-editor")
	if got := DetectEditor(); got != "qc-visual-editor" {
		t.Errorf("DetectEditor() = %q, want VISUAL to win", got)
	}

	t.Setenv("VISUAL", "")
	if got := DetectEditor(); got != "qc-env-editor" {
		t.Errorf("DetectEditor() = %q, want EDITOR to be used", got)
	}

	// An editor that is not installed must not be returned.
	t.Setenv("EDITOR", "qc-missing-editor")
	if got := DetectEditor(); got != "vim" {
		t.Errorf("DetectEditor() = %q, want the probed fallback %q", got, "vim")
	}
}

func TestDetectEditorWithoutAnyEditor(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "")

	if got := DetectEditor(); got != "" {
		t.Errorf("DetectEditor() = %q, want an empty string", got)
	}
}

func TestResolveEditor(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell based fake editors are not portable to Windows")
	}
	dir := t.TempDir()
	fakeEditor(t, dir, "qc-configured-editor")
	fakeEditor(t, dir, "vim")
	t.Setenv("PATH", dir)
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "")

	editor, replaced, err := ResolveEditor("qc-configured-editor")
	if err != nil || replaced || editor != "qc-configured-editor" {
		t.Errorf(
			"ResolveEditor() = (%q, %v, %v), want the configured editor",
			editor,
			replaced,
			err,
		)
	}

	// A configured editor that is not installed falls back to a detected one.
	editor, replaced, err = ResolveEditor("qc-missing-editor")
	if err != nil || !replaced || editor != "vim" {
		t.Errorf(
			"ResolveEditor() = (%q, %v, %v), want the fallback %q",
			editor,
			replaced,
			err,
			"vim",
		)
	}

	// An empty config entry is resolved without reporting a replacement.
	editor, replaced, err = ResolveEditor("")
	if err != nil || replaced || editor != "vim" {
		t.Errorf("ResolveEditor() = (%q, %v, %v), want the detected editor", editor, replaced, err)
	}

	t.Setenv("PATH", t.TempDir())
	if _, _, err = ResolveEditor("qc-missing-editor"); err != ErrNoEditor {
		t.Errorf("ResolveEditor() error = %v, want ErrNoEditor", err)
	}
}
