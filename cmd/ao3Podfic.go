package cmd

import (
	"log"

	"github.com/ohzqq/urmeta/ao3"
	"github.com/spf13/cobra"
)

// ao3PodficCmd represents the ao3Podfic command
var ao3PodficCmd = &cobra.Command{
	Use:     "podfic",
	Aliases: []string{"p"},
	Short:   "scrape a podfic",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s, err := ao3.Work(args[0], true)
		if err != nil {
			log.Fatal(err)
		}
		processMetadata(s)
	},
}

func init() {
	ao3Cmd.AddCommand(ao3PodficCmd)
}
