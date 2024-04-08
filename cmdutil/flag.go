package cmdutil

import (
	"fmt"

	"github.com/spf13/cobra"
)

/* FUNCTION */
func addFlag(cmd *cobra.Command, flag Flag) {
	switch flag.Type {
	case "string":
		{
			cmd.Flags().String(flag.Flag, flag.Default, flag.Description)
		}
	case "bool":
		{
			cmd.Flags().Bool(flag.Flag, flag.Default == "true", flag.Description)
		}
	default:
		{
			fmt.Println("Unknown flag type:", flag.Type)
		}
	}
}

func AddEditFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("edit", "e", false, "Edit the file where the command is defined")
}

func AddFlags(cmd *cobra.Command, flags []Flag) {
	for _, flag := range flags {
		addFlag(cmd, flag)
	}
}
