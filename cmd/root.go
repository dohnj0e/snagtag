package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "snagtag",
	Short: "SnagTag - cli tool for scraping social media platforms",
	Long:  "SnagTag - cli tool for scraping social media platforms, such as Youtube and Tiktok",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: sub-command is required. Run 'snagtag --help' for usage information.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
