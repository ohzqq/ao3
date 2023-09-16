package cmd

import (
	"log"

	"github.com/ohzqq/urmeta/ao3"
	"github.com/spf13/cobra"
)

// ao3WorkCmd represents the ao3Work command
var ao3WorkCmd = &cobra.Command{
	Use:     "work",
	Aliases: []string{"w"},
	Short:   "scrape a work",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s, err := ao3.Work(args[0], isPodfic)
		if err != nil {
			log.Fatal(err)
		}
		processMetadata(s)
	},
}

func init() {
	ao3Cmd.AddCommand(ao3WorkCmd)
}
