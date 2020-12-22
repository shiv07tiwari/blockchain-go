package main

import (
	"github.com/spf13/cobra"
)

// A simple command to print the version.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Describes version.",
	Run: func(cmd *cobra.Command, args []string) {
		print("A version command for the CLI\n")
	},
}
