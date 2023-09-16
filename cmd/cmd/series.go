package cmd

import (
	"log"

	"github.com/ohzqq/ao3"
	"github.com/spf13/cobra"
)

// seriesCmd represents the series command
var seriesCmd = &cobra.Command{
	Use:     "series",
	Aliases: []string{"s"},
	Short:   "scrape a series",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s, err := ao3.Page(args[0], isPodfic)
		if err != nil {
			log.Fatal(err)
		}
		processMetadata(s)
	},
}

func init() {
	rootCmd.AddCommand(seriesCmd)
}
