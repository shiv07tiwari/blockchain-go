// Child command to print the balance sheet for the cafe, for all the users
package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var balancesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all the User balances.",
	Run: func(cmd *cobra.Command, args []string) {

		// Get the state variable of the blockchain.
		state, err := BlockChain.GetState()

		// Handle the error and log the required data
		if err != nil {
			fmt.Println(err)
		}

		// Close the db Instance using defer, i.e. once the function returns to the caller.
		defer state.Close()

		fmt.Println("Balance Sheet")
		fmt.Println("__________________")

		// Iterate over all the balances in our state, and print out the account name and balance.
		for account, balance := range state.Balances {
			fmt.Println(fmt.Sprintf("%s: %d", account, balance))
		}
	},
}
