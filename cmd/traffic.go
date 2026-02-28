package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mrm/keywords-everywhere-cli/internal/api"
	"github.com/spf13/cobra"
)

var domainKeywordsCmd = &cobra.Command{
	Use:   "domain-keywords <domain>",
	Short: "Get keywords a domain ranks for",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetDomainKeywords(args[0], country, currency, limit)
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

var urlKeywordsCmd = &cobra.Command{
	Use:   "url-keywords <url>",
	Short: "Get keywords a URL ranks for",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetURLKeywords(args[0], country, currency, limit)
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

var domainTrafficCmd = &cobra.Command{
	Use:   "domain-traffic <domain>",
	Short: "Get traffic metrics for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetDomainTraffic(args[0])
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

var urlTrafficCmd = &cobra.Command{
	Use:   "url-traffic <url>",
	Short: "Get traffic metrics for a URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetURLTraffic(args[0])
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
	rootCmd.AddCommand(domainKeywordsCmd)
	rootCmd.AddCommand(urlKeywordsCmd)
	rootCmd.AddCommand(domainTrafficCmd)
	rootCmd.AddCommand(urlTrafficCmd)
}
