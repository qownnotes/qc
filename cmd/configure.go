package cmd

import (
	"github.com/qownnotes/qc/config"
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Edit config file",
	Long:  `Edit config file (default: opened by vim)`,
	RunE:  configure,
}

func configure(cmd *cobra.Command, args []string) (err error) {
	editor := config.Conf.General.Editor
	return editFile(editor, configFile)
}

func init() {
	RootCmd.AddCommand(configureCmd)
}
