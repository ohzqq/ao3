package cmd

import (
	"log"

	"github.com/ohzqq/urmeta/ao3"
	"github.com/spf13/cobra"
)

// ao3SeriesCmd represents the ao3Series command
var ao3SeriesCmd = &cobra.Command{
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
	ao3Cmd.AddCommand(ao3SeriesCmd)
}
