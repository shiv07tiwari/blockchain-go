package main

import (
	"blockchain-go/data"
	"fmt"

	"github.com/spf13/cobra"
)

// CLI command to print out the Balances of all the Users.
var balencesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all the User balances.",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := data.NewStateFromDisk()
		if err != nil {
			fmt.Println(err)
		}
		defer state.Close()
		fmt.Println("Balance Sheet")
		fmt.Println("__________________")
		fmt.Println("")

		for account, balance := range state.Balances {
			fmt.Println(fmt.Sprintf("%s: %d", account, balance))
		}
	},
}
