package cmd

import (
	"os"

	"github.com/ohzqq/urmeta"
	"github.com/spf13/cobra"
)

var query = &urmeta.Query{}
var uriArg string
var noSave bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "urmeta",
	Short: "scrape book metadata",
	Long:  `scrape book/work metadata from audible and ao3`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&uriArg, "url", "u", "", "scrape url")
	rootCmd.PersistentFlags().BoolVar(&noSave, "no-save", false, "don't write metadata to disk")
}
