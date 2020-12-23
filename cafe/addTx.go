// Child command to add a new transaction in the blockchain
package main

import (
	"blockchain-go/data"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func getAddTxCommand() *cobra.Command {

	var addTxCommand = &cobra.Command{
		Use:   "add",
		Short: "Adds a new Transaction to Database",
		Run: func(cmd *cobra.Command, args []string) {

			// Get the values entered for the respective flags
			from, _ := cmd.Flags().GetString("from")
			to, _ := cmd.Flags().GetString("to")
			amount, _ := cmd.Flags().GetInt("amount")
			txData, _ := cmd.Flags().GetString("data")

			// Get the state variable of the blockchain.
			state, err := BlockChain.GetState()

			// Create the transaction
			tx, err := data.NewTransaction(from, to, amount, state)

			// Handle the error and log the required data
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				fmt.Println("Tx detail : ", txData)
				os.Exit(1)
			}

			// Close the db Instance using defer, i.e. once the function returns to the caller.
			defer state.Close()

			// Add the new transaction to the state
			err = state.Add(tx)

			// Handle the error and log the required data
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			// Write the new transaction to the persistent database.
			_, err = state.Persist()

			// Handle the error and log the required data to standard error file descriptor
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			// Print out the success message
			fmt.Println("Transaction Successful")

		},
	}

	// Add flags to the command and mark them as required.(as these are required values for a new Tx)
	addTxCommand.Flags().String("from", "null", "Paid From")
	addTxCommand.MarkFlagRequired("from")
	addTxCommand.Flags().String("to", "null", "Paid To")
	addTxCommand.MarkFlagRequired("to")
	addTxCommand.Flags().Int("amount", 0, "Paid Amount")
	addTxCommand.MarkFlagRequired("amount")

	// Add description flag, which is optional.
	addTxCommand.Flags().String("data", "random", "Details")

	return addTxCommand
}
