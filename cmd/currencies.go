package cmd

import (
	"fmt"

	"github.com/mrm/keywords-everywhere-cli/internal/api"
	"github.com/spf13/cobra"
)

var currenciesCmd = &cobra.Command{
	Use:   "currencies",
	Short: "List supported currencies",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetCurrencies()
		if err != nil {
			return err
		}

		fmt.Println(string(result))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(currenciesCmd)
}
