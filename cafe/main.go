package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// BlockChain state variable
var BlockChain *Blockchain

func main() {
	var tbbCmd = &cobra.Command{
		Use:   "cafe",
		Short: "The Blockchain Cafe CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	BlockChain, _ = NewBlockChain()
	tbbCmd.AddCommand(versionCmd)
	tbbCmd.AddCommand(balencesListCmd)
	tbbCmd.AddCommand(getAddTxCommand())

	err := tbbCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
