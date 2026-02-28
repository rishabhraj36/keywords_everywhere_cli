package cmd

import (
	"fmt"

	"github.com/mrm/keywords-everywhere-cli/internal/api"
	"github.com/spf13/cobra"
)

var countriesCmd = &cobra.Command{
	Use:   "countries",
	Short: "List supported countries",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetCountries()
		if err != nil {
			return err
		}

		fmt.Println(string(result))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(countriesCmd)
}
