package cmd

import (
	"errors"
	"io"
	"os/exec"
	"runtime"
	"testing"
)

func TestRunExitStatus(t *testing.T) {
	for _, tt := range []struct {
		name    string
		command string
		code    int
	}{
		{"success", "exit 0", 0},
		{"interrupted", "exit 130", 130},
		{"failure", "exit 7", 7},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := run(tt.command, nil, io.Discard)
			switch tt.code {
			case 0:
				if err != nil {
					t.Fatal(err)
				}
			case 130:
				if !errors.Is(err, ErrSigInt) {
					t.Fatalf("expected interruption, got %v", err)
				}
			default:
				var exitErr *exec.ExitError
				if !errors.As(err, &exitErr) || exitErr.ExitCode() != tt.code {
					t.Fatalf("expected exit code %d, got %v", tt.code, err)
				}
			}
		})
	}
}

func TestRunSignal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("requires POSIX signals")
	}
	for _, tt := range []struct {
		name        string
		command     string
		interrupted bool
	}{
		{"SIGINT", "kill -INT $$", true},
		{"SIGTERM", "kill -TERM $$", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := run(tt.command, nil, io.Discard)
			if err == nil || errors.Is(err, ErrSigInt) != tt.interrupted {
				t.Fatalf("expected interrupted=%v, got %v", tt.interrupted, err)
			}
		})
	}
}
