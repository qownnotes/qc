package cmd

import (
	"bufio"
	"fmt"
	"github.com/qownnotes/qc/config"
	"github.com/qownnotes/qc/websocket"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch note folder",
	Long:  `Switch to a different note folder`,
	RunE:  switchNoteFolder,
}

func switchNoteFolder(cmd *cobra.Command, args []string) (err error) {
	flag := config.Flag

	if flag.Query != "" {
		id, _ := strconv.Atoi(flag.Query)
		fmt.Printf("Attempting to switch to note folder number %d!\n", id)
		websocket.SwitchNoteFolder(id)

		return nil
	}

	noteFolderData, currentId := websocket.FetchNoteFolderData()

	for _, noteFolder := range noteFolderData {
		currentText := ""

		if currentId == noteFolder.Id {
			currentText = "*"
		}

		fmt.Printf("%d%s) %s\n", noteFolder.Id, currentText, noteFolder.Name)
	}

	fmt.Print("\nSelect note folder to switch to: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading standard input:", err)

		return nil
	}

	id, _ := strconv.Atoi(scanner.Text())

	for _, noteFolder := range noteFolderData {
		if noteFolder.Id == id {
			websocket.SwitchNoteFolder(id)
			return nil
		}
	}

	fmt.Printf("Could not find note folder number %d!\n\n", id)

	return switchNoteFolder(cmd, args)
}

func init() {
	RootCmd.AddCommand(switchCmd)
	switchCmd.Flags().StringVarP(&config.Flag.Query, "folder", "f", "",
		`Note folder id to switch to`)
}
