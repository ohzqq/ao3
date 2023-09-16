package cmd

import (
	"log"

	"github.com/ohzqq/ao3"
	"github.com/spf13/cobra"
)

// podficCmd represents the podfic command
var podficCmd = &cobra.Command{
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
	rootCmd.AddCommand(podficCmd)
}
