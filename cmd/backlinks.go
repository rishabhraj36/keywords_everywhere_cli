package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mrm/keywords-everywhere-cli/internal/api"
	"github.com/spf13/cobra"
)

var domainBacklinksCmd = &cobra.Command{
	Use:   "domain-backlinks <domain>",
	Short: "Get backlinks for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetDomainBacklinks(args[0], limit)
		if err != nil {
			return err
		}

		output, err := json.Marshal(result)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

var pageBacklinksCmd = &cobra.Command{
	Use:   "page-backlinks <url>",
	Short: "Get backlinks for a specific page",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetPageBacklinks(args[0], limit)
		if err != nil {
			return err
		}

		output, err := json.Marshal(result)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(domainBacklinksCmd)
	rootCmd.AddCommand(pageBacklinksCmd)
}
