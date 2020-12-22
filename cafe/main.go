package main

import (
	"fmt"
	"os"

	// Cobra is used to generate the CLI
	"github.com/spf13/cobra"
)

// BlockChain state variable
var BlockChain *Blockchain

// This is the entry point of the application
func main() {

	// Root command for The Blockchain Cafe
	var tbcCmd = &cobra.Command{
		Use:   "cafe",
		Short: "The Blockchain Cafe CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// Get a new blockchain instance
	BlockChain, _ = NewBlockChain()

	// add child command to get the version
	tbcCmd.AddCommand(versionCmd)

	// add child command to get the balances of every user
	tbcCmd.AddCommand(balancesListCmd)

	// add child command to create a new transaction
	tbcCmd.AddCommand(getAddTxCommand())

	// Execute the root command and handle the error if any.
	err := tbcCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
