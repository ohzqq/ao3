package cmd

import (
	"fmt"
	"strings"

	"github.com/ohzqq/cdb"
	"github.com/ohzqq/urmeta"
	"github.com/ohzqq/urmeta/ao3"
	"github.com/spf13/cobra"
)

var (
	isPodfic    bool
	formats     []string
	noDownloads bool
)

// ao3Cmd represents the ao3 command
var ao3Cmd = &cobra.Command{
	Use:   "ao3",
	Short: "scrape ao3",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func processMetadata(books []cdb.Book) {
	for _, b := range books {
		fmt.Printf("%#v\n", b)
		if !noSave {
			urmeta.SaveMetadata(b)
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

func init() {
	rootCmd.AddCommand(ao3Cmd)
	ao3Cmd.PersistentFlags().StringSliceVarP(&formats, "formats", "f", []string{".epub"}, "format to download")
	ao3Cmd.PersistentFlags().BoolVarP(&noDownloads, "no-downloads", "d", false, "don't download any formats")
	ao3Cmd.MarkFlagsMutuallyExclusive("formats", "no-downloads")

	ao3Cmd.PersistentFlags().BoolVarP(&isPodfic, "podfic", "p", false, "scrape podfic url")
}
