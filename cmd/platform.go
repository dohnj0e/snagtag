package cmd

import (
	"fmt"

	"github.com/dohnj0e/snagtag/platforms/tiktok"
	"github.com/dohnj0e/snagtag/platforms/twitter"
	"github.com/dohnj0e/snagtag/platforms/youtube"
	"github.com/spf13/cobra"
)

var keyword string

var platformCmd = &cobra.Command{
	Use:   "platform",
	Short: "Scrape a specific platform for a given keyword",
}

var youtubeCmd = &cobra.Command{
	Use:   "youtube",
	Short: "Scrape youtube for a given keyword",
	Run: func(cmd *cobra.Command, args []string) {
		if keyword == "" {
			fmt.Println("Error: no key provided. Use '--keyword (-k)' to specify a keyword.")
			return
		}

		err := youtube.Scrape(keyword)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	},
}

var tiktokCmd = &cobra.Command{
	Use:   "tiktok",
	Short: "Scrape tiktok for a given keyword",
	Run: func(cmd *cobra.Command, args []string) {
		if keyword == "" {
			fmt.Println("Error: no key provided. Use '--keyword (-k)' to specify a keyword.")
			return
		}

		err := tiktok.Scrape(keyword)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	},
}

var twitterCmd = &cobra.Command{
	Use:   "twitter",
	Short: "Scrape twitter for a given keyword",
	Run: func(cmd *cobra.Command, args []string) {
		if keyword == "" {
			fmt.Println("Error: no keyword provided. Use '--keyword (-k)' to specify a keyword.")
			return
		}

		err := twitter.Scrape(keyword)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(platformCmd)
	platformCmd.AddCommand(youtubeCmd, tiktokCmd, twitterCmd)
	platformCmd.PersistentFlags().StringVarP(&keyword, "keyword", "k", "", "Keyword to search for")
}
