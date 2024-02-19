package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/danielgtaylor/casing"
	"github.com/ohzqq/ao3"
	"github.com/ohzqq/cdb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	isPodfic    bool
	formats     []string
	noDownloads bool
	query       = &ao3.Query{}
	flagURI     string
	noSave      bool
	flagEncode  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ao3",
	Short: "scrape ao3 metadata",
	Long:  `scrape book/work metadata from ao3`,
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&flagURI, "url", "u", "", "scrape url")
	rootCmd.PersistentFlags().BoolVar(&noSave, "no-save", false, "don't write metadata to disk")
	rootCmd.PersistentFlags().BoolVarP(&isPodfic, "podfic", "p", false, "scrape podfic url")

	rootCmd.PersistentFlags().StringVarP(&flagEncode, "encode", "e", ".yaml", "encode [.yaml|.toml|.json]")

	rootCmd.PersistentFlags().StringSliceVarP(&formats, "formats", "f", []string{".epub"}, "format to download")
	rootCmd.PersistentFlags().BoolVarP(&noDownloads, "no-downloads", "d", false, "don't download any formats")
	rootCmd.MarkFlagsMutuallyExclusive("formats", "no-downloads")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".toot" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".toot")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func processMetadata(books []cdb.Book) {
	for _, b := range books {
		if !noSave {
			name := casing.Snake(b.Title) + flagEncode
			err := b.Save(name, true)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := b.Print(flagEncode, true)
			if err != nil {
				log.Fatal(err)
			}
		}
		if !noDownloads {
			downloadFormats(b)
		}
	}
}

func downloadFormats(b cdb.Book) {
	for _, f := range b.Formats {
		for _, ext := range formats {
			if strings.Contains(f, ext) {
				fmt.Printf("downloading %s\n", b.Title+ext)
				name := casing.Snake(b.Title) + ext
				ao3.DownloadWork(f, name)
			}
		}
	}
}
