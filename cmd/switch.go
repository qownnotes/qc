package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/qownnotes/qc/config"
	"github.com/qownnotes/qc/entity"
	"github.com/qownnotes/qc/websocket"
	"github.com/spf13/cobra"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch note folder",
	Long:  `Choose a note folder with the configured selector, or switch directly with --folder.`,
	RunE:  switchNoteFolder,
}

func switchNoteFolder(cmd *cobra.Command, args []string) (err error) {
	flag := config.Flag

	if flag.Query != "" {
		id, err := strconv.Atoi(flag.Query)
		if err != nil || id < 0 {
			return fmt.Errorf("invalid note folder id: %q", flag.Query)
		}
		fmt.Printf("Attempting to switch to note folder number %d!\n", id)
		websocket.SwitchNoteFolder(id)

		return nil
	}

	noteFolderData, currentId := websocket.FetchNoteFolderData()
	id, selected, err := selectNoteFolder(noteFolderData, currentId)
	if err != nil {
		return err
	}
	if selected {
		websocket.SwitchNoteFolder(id)
	}
	return nil
}

func selectNoteFolder(folders []entity.NoteFolderInfo, currentId int) (int, bool, error) {
	if len(folders) == 0 {
		return 0, false, fmt.Errorf("no note folders available")
	}
	command := strings.TrimSpace(config.Conf.General.SelectCmd)
	if command == "" {
		return 0, false, fmt.Errorf("no selector configured; set general.selectcmd in qc configure")
	}
	// Keep custom selectors usable, but give fzf a compact, single-choice UI.
	selector := filepath.Base(strings.Trim(strings.Fields(command)[0], `"'`))
	if selector == "fzf" || selector == "fzf.exe" {
		command += ` --height=60% --layout=reverse --border=rounded --no-multi` +
			` --delimiter="\t" --with-nth=2.. --prompt="Note folder > "` +
			` --header="Enter: switch | Esc: cancel | Type to search"`
	}

	var input strings.Builder
	choices := make(map[string]int, len(folders))
	cleanName := strings.NewReplacer("\n", " ", "\r", " ", "\t", " ")
	// Put the current folder first without changing the server's remaining order.
	for _, current := range []bool{true, false} {
		for _, folder := range folders {
			if (folder.Id == currentId) != current {
				continue
			}
			label := cleanName.Replace(folder.Name)
			if current {
				label += " (current)"
			}
			line := fmt.Sprintf("%d\t%s", folder.Id, label)
			choices[line] = folder.Id
			input.WriteString(line + "\n")
		}
	}

	var output bytes.Buffer
	if err := run(command, strings.NewReader(input.String()), &output); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			// fzf/peco use status 1 when no item was selected.
			return 0, false, nil
		}
		return 0, false, err
	}
	line := strings.TrimSuffix(strings.TrimSuffix(output.String(), "\n"), "\r")
	if line == "" {
		return 0, false, nil
	}
	id, ok := choices[line]
	if !ok {
		return 0, false, fmt.Errorf("selector returned an unknown note folder")
	}
	return id, true, nil
}

func init() {
	RootCmd.AddCommand(switchCmd)
	switchCmd.Flags().StringVarP(&config.Flag.Query, "folder", "f", "",
		`Note folder id to switch to`)
}
