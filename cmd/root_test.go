package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/qownnotes/qc/config"
	"github.com/spf13/cobra"
)

func TestExecuteHelper(t *testing.T) {
	mode := os.Getenv("QC_EXECUTE_TEST")
	if mode == "" {
		return
	}
	if strings.HasPrefix(mode, "search-") || strings.HasPrefix(mode, "exec-") {
		// Use cached snippets and a closed socket port so no QOwnNotes instance
		// or user configuration is needed by the subprocess test.
		dir, err := config.GetDefaultConfigDir()
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "snippets.cache"),
			[]byte(`{"type":"commandSnippets","data":[]}`), 0600); err != nil {
			t.Fatal(err)
		}
		config.Conf.QOwnNotes.WebSocketPort = 0
		config.Conf.General.SelectCmd = "exit 130"
		if strings.HasSuffix(mode, "signal") {
			config.Conf.General.SelectCmd = "kill -INT $$"
		} else if strings.HasSuffix(mode, "cancel") {
			config.Conf.General.SelectCmd = "exit 1"
		}
		// Cobra's registered initializer reloads the config before RunE.
		configFile = filepath.Join(dir, "config.toml")
		contents := fmt.Sprintf("[General]\nselectcmd = %q\n[QOwnNotes]\nwebsocket_port = 0\n",
			config.Conf.General.SelectCmd)
		if err := os.WriteFile(configFile, []byte(contents), 0600); err != nil {
			t.Fatal(err)
		}
		log.SetOutput(io.Discard)
	}
	RootCmd = &cobra.Command{
		Use: "qc",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.HasPrefix(mode, "search-") {
				return search(cmd, args)
			}
			if strings.HasPrefix(mode, "exec-") {
				return execute(cmd, args)
			}
			if mode == "signal" {
				fmt.Println("ready")
				select {}
			}
			return fmt.Errorf("wrapped: %w", ErrSigInt)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	RootCmd.SetArgs(nil)
	Execute()
	os.Exit(0)
}

func TestExecuteInterrupt(t *testing.T) {
	for _, mode := range []string{
		"error", "signal", "search-interrupt", "exec-interrupt",
		"search-signal", "exec-signal", "search-cancel",
	} {
		t.Run(mode, func(t *testing.T) {
			if strings.HasSuffix(mode, "signal") && runtime.GOOS == "windows" {
				t.Skip("os.Process.Signal does not support interrupt on Windows")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestExecuteHelper$")
			home := t.TempDir()
			command.Env = append(os.Environ(), "QC_EXECUTE_TEST="+mode,
				"HOME="+home, "APPDATA="+home, "USERPROFILE="+home)
			var output []byte
			var err error
			if mode == "signal" {
				var stderr bytes.Buffer
				command.Stderr = &stderr
				stdout, pipeErr := command.StdoutPipe()
				if pipeErr != nil {
					t.Fatal(pipeErr)
				}
				if err = command.Start(); err != nil {
					t.Fatal(err)
				}
				scanner := bufio.NewScanner(stdout)
				if !scanner.Scan() || scanner.Text() != "ready" {
					_ = command.Wait()
					t.Fatal("helper did not become ready")
				}
				if err = command.Process.Signal(os.Interrupt); err != nil {
					t.Fatal(err)
				}
				for scanner.Scan() {
					output = append(output, scanner.Bytes()...)
				}
				err = command.Wait()
				output = append(output, stderr.Bytes()...)
			} else {
				output, err = command.CombinedOutput()
			}
			if mode == "search-cancel" {
				if err != nil {
					t.Fatalf(
						"expected ordinary selector cancellation to succeed, got %v; output: %s",
						err,
						output,
					)
				}
			} else {
				exitErr, ok := err.(*exec.ExitError)
				if !ok || exitErr.ExitCode() != 130 {
					t.Fatalf("expected exit code 130, got %v; output: %s", err, output)
				}
			}
			if len(output) != 0 {
				t.Fatalf("expected quiet cancellation, got %q", output)
			}
		})
	}
}
