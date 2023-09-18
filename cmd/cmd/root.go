package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ohzqq/ao3"
	"github.com/ohzqq/cdb"
	"github.com/spf13/cobra"
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
	rootCmd.PersistentFlags().StringVarP(&flagURI, "url", "u", "", "scrape url")
	rootCmd.PersistentFlags().BoolVar(&noSave, "no-save", false, "don't write metadata to disk")
	rootCmd.PersistentFlags().BoolVarP(&isPodfic, "podfic", "p", false, "scrape podfic url")

	rootCmd.PersistentFlags().StringVarP(&flagEncode, "encode", "e", ".yaml", "encode [.yaml|.toml|.json]")

	rootCmd.PersistentFlags().StringSliceVarP(&formats, "formats", "f", []string{".epub"}, "format to download")
	rootCmd.PersistentFlags().BoolVarP(&noDownloads, "no-downloads", "d", false, "don't download any formats")
	rootCmd.MarkFlagsMutuallyExclusive("formats", "no-downloads")
}

func processMetadata(books []cdb.Book) {
	for _, b := range books {
		//fmt.Printf("%#v\n", b)
		if !noSave {
			err := b.Save(b.Title, flagEncode, editable)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := b.Print(flagEncode, editable)
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
				ao3.DownloadWork(f, b.Title+ext)
			}
		}
	}
}
