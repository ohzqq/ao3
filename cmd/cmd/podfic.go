package cmd

import (
	"log"

	"github.com/ohzqq/ao3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// podficCmd represents the podfic command
var podficCmd = &cobra.Command{
	Use:     "podfic",
	Aliases: []string{"p"},
	Short:   "scrape a podfic",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("podfic", true)
		viper.Set("no-downloads", true)

		s, err := ao3.Scrape(args[0])
		if err != nil {
			log.Fatal(err)
		}
		processMetadata(s)
	},
}

func init() {
	rootCmd.AddCommand(podficCmd)
}
