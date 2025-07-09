package commands

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "peek",
	Short: "Scan and interact with nearby Bluetooth devices.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("no command specified")
	},
}

// Execute runs the root command and handles subcommands.
func Execute() error {
	rootCommand.SetOut(os.Stdout)
	return rootCommand.Execute()
}
