package cmd

import (
	"github.com/spf13/cobra"
)

var sowCmd = &cobra.Command{
	Use:   "sow",
	Short: "scraping subcommand",
	Long:  "use this module to scrape pricing data",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(sowCmd)

	addDBFlags(sowCmd)
}
