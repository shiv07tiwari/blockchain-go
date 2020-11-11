package main

import (
	"blockchain-go/data"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func getAddTxCommand() *cobra.Command {
	// command to add the transaction
	var addTxCommand = &cobra.Command{
		Use:   "add",
		Short: "Adds a new Transaction to Database",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString("from")
			to, _ := cmd.Flags().GetString("to")
			amount, _ := cmd.Flags().GetUint("amount")
			txData, _ := cmd.Flags().GetString("data")

			tx := data.NewTx(from, to, amount, txData)
			state, err := BlockChain.GetState()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()
			err = state.Add(tx)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			_, err = state.Persist()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Println("Transaction Successful")

		},
	}

	addTxCommand.Flags().String("from", "null", "Paid From")
	addTxCommand.MarkFlagRequired("from")
	addTxCommand.Flags().String("to", "null", "Paid To")
	addTxCommand.MarkFlagRequired("to")
	addTxCommand.Flags().Uint("amount", 0, "Paid Amount")
	addTxCommand.MarkFlagRequired("amount")
	addTxCommand.Flags().String("data", "random", "Details")

	return addTxCommand
}
