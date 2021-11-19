package cmd

import (
	"fmt"
	"os"

	"github.com/qownnotes/qc/config"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:           "qc",
	Short:         "QOwnNotes command-line snippet manager.",
	Long:          `qc - QOwnNotes command-line snippet manager.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	// cobra.OnInitialize(initConfig)
	RootCmd.AddCommand(versionCmd)

	RootCmd.PersistentFlags().IntVarP(&config.Flag.Port, "port", "p", 22322, "socket server port in QOwnNotes")
	RootCmd.PersistentFlags().BoolVarP(&config.Flag.Debug, "debug", "", false, "debug mode")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("qc version %s\n", version)
	},
}
