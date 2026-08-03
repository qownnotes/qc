package cmd

import (
	"fmt"
	"os"

	"github.com/qownnotes/qc/config"
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Edit config file",
	Long: `Edit config file.

The editor is taken from the "editor" option of the config file. If that
editor is not set or not installed, the VISUAL and EDITOR environment
variables are used, and finally a list of common editors is probed.`,
	RunE: configure,
}

func configure(cmd *cobra.Command, args []string) (err error) {
	editor, replaced, err := config.ResolveEditor(config.Conf.General.Editor)
	if err != nil {
		return err
	}
	if replaced {
		fmt.Fprintf(os.Stderr,
			"qc: the configured editor %q is not available, using %q instead\n",
			config.Conf.General.Editor, editor)
	}
	return editFile(editor, configFile)
}

func init() {
	RootCmd.AddCommand(configureCmd)
}
