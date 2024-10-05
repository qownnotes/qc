package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/qownnotes/qc/config"
	"github.com/spf13/cobra"
	"gopkg.in/alessio/shellescape.v1"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Run the selected commands",
	Long:  `Run the selected commands directly`,
	RunE:  execute,
}

var (
	lastCmdFile string
)

func execute(cmd *cobra.Command, args []string) (err error) {
	flag := config.Flag

	var options []string
	var command string
	var writeLastCmd bool

	if flag.Query != "" {
		options = append(options, fmt.Sprintf("--query %s", shellescape.Quote(flag.Query)))
	}

	if config.Flag.Last {
		command = readLastCmdFile()
	}

	if command == "" {
		commands, err := filter(options, flag.FilterTag)
		if err != nil {
			return err
		}
		command = strings.Join(commands, "; ")
		writeLastCmd = true
	}

	if config.Flag.Debug {
		fmt.Printf("Command: %s\n", command)
	}

	if config.Flag.Command {
		fmt.Printf("%s: %s\n", color.YellowString("Command"), command)
	}

	if command == "" {
		return nil
	}

	if writeLastCmd {
		// store last command
		writeLastCmdFile(command)
	}

	// Add commands to the Atuin history
	if config.Flag.Atuin {
		escapedCommand := escapeCommandForShell(command)
		command = `histid=$(atuin history start -- ` + escapedCommand + `` + ")\n" + command +
			"\natuin history end --exit $? $histid"
	}

	return run(command, os.Stdin, os.Stdout)
}

// escapeCommandForShell escapes a command for use in a shell script.
// It also works for multi-line commands.
func escapeCommandForShell(command string) string {
	// Trim and split the command into lines
	lines := strings.Split(strings.Trim(command, "\n"), "\n")

	// Quote each line individually and join with literal $'\n'
	quotedLines := make([]string, len(lines))
	for i, line := range lines {
		// Use single quotes to avoid most escaping issues
		quotedLines[i] = "'" + strings.ReplaceAll(line, "'", "'\"'\"'") + "'"
	}

	return strings.Join(quotedLines, "$'\\n'")
}

func init() {
	RootCmd.AddCommand(execCmd)
	execCmd.Flags().BoolVarP(&config.Flag.Color, "color", "", false,
		`Enable colorized output (only fzf)`)
	execCmd.Flags().StringVarP(&config.Flag.Query, "query", "q", "",
		`Initial value for query`)
	execCmd.Flags().StringVarP(&config.Flag.FilterTag, "tag", "t", "",
		`Filter tag`)
	execCmd.Flags().BoolVarP(&config.Flag.Command, "command", "c", false,
		`Show the command with the plain text before executing`)
	execCmd.Flags().BoolVarP(&config.Flag.Last, "last", "l", false,
		`Execute the last command`)
	execCmd.Flags().BoolVarP(&config.Flag.Atuin, "atuin", "a", false,
		`Store single-line command in Atuin history`)

	initLastCmdFile()
}

func initLastCmdFile() {
	if lastCmdFile == "" {
		dir, err := config.GetDefaultConfigDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}

		lastCmdFile = filepath.Join(dir, "lastcmd")
	}
}

func writeLastCmdFile(cmd string) {
	if err := os.WriteFile(lastCmdFile, []byte(cmd), 0600); err != nil {
		log.Fatal("Could not write last command file: ", err)
	}
}

func readLastCmdFile() string {
	_, err := os.Stat(lastCmdFile)

	if errors.Is(err, os.ErrNotExist) {
		return ""
	}

	data, err := os.ReadFile(lastCmdFile)

	if err != nil {
		log.Fatal("Could not read last command file: ", err)
	}

	return string(data)
}
