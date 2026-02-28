package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	country  string
	currency string
	source   string
	limit    int
)

var rootCmd = &cobra.Command{
	Use:   "ke",
	Short: "Keywords Everywhere CLI",
	Long:  "A CLI for the Keywords Everywhere API. Set KEYWORDS_EVERYWHERE_API_KEY environment variable.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&country, "country", "c", "us", "Country code")
	rootCmd.PersistentFlags().StringVar(&currency, "currency", "usd", "Currency code")
	rootCmd.PersistentFlags().StringVarP(&source, "source", "s", "gkp", "Data source: gkp|cli")
	rootCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 0, "Max results (0 = no limit)")
}
