package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mrm/keywords-everywhere-cli/internal/api"
	"github.com/spf13/cobra"
)

var keywordsCmd = &cobra.Command{
	Use:   "keywords [keyword...]",
	Short: "Get volume, CPC, and competition data for keywords",
	Long:  "Get keyword data. Pass keywords as args or pipe via stdin (one per line).",
	RunE: func(cmd *cobra.Command, args []string) error {
		keywords := args

		// Read from stdin if no args and stdin has data
		if len(keywords) == 0 {
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					line := scanner.Text()
					if line != "" {
						keywords = append(keywords, line)
					}
				}
			}
		}

		if len(keywords) == 0 {
			return fmt.Errorf("no keywords provided")
		}

		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.GetKeywordData(keywords, country, currency, source)
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
	rootCmd.AddCommand(keywordsCmd)
}
