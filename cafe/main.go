package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var tbbCmd = &cobra.Command{
		Use:   "cafe",
		Short: "The Blockchain Cafe CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	tbbCmd.AddCommand(versionCmd)
	tbbCmd.AddCommand(balencesListCmd)

	err := tbbCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
