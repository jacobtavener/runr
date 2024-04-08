package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

/* GLOBALS */
var rootCmd = &cobra.Command{
	Use:   "runr",
	Short: "A CLI tool for discovering and running scripts in a project",
	Long:  ``,
}

/* FUNCTIONS */
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
