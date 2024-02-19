package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/danielgtaylor/casing"
	"github.com/ohzqq/ao3"
	"github.com/ohzqq/audbk"
	"github.com/ohzqq/cdb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	query  = &ao3.Query{}
	ffmeta bool
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

	rootCmd.PersistentFlags().StringP("url", "u", "", "scrape url")
	rootCmd.PersistentFlags().Bool("no-save", false, "don't write metadata to disk")

	rootCmd.PersistentFlags().BoolP("podfic", "p", false, "scrape podfic url")
	rootCmd.PersistentFlags().BoolVarP(&ffmeta, "ffmeta", "m", false, "write ffmeta")

	rootCmd.PersistentFlags().StringP("encode", "e", ".yaml", "encode [.yaml|.toml|.json|.ini]")

	rootCmd.PersistentFlags().StringSliceP("formats", "f", []string{".epub"}, "format to download")
	rootCmd.PersistentFlags().BoolP("no-downloads", "d", false, "don't download any formats")
	rootCmd.MarkFlagsMutuallyExclusive("formats", "no-downloads")

	viper.BindPFlag("no-save", rootCmd.PersistentFlags().Lookup("no-save"))
	viper.BindPFlag("no-downloads", rootCmd.PersistentFlags().Lookup("no-downloads"))
	viper.BindPFlag("podfics", rootCmd.PersistentFlags().Lookup("podfics"))
	viper.BindPFlag("formats", rootCmd.PersistentFlags().Lookup("formats"))
	viper.BindPFlag("encode", rootCmd.PersistentFlags().Lookup("encode"))
}

func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	viper.SetDefault("podfic", false)
	viper.SetDefault("no-save", false)
	viper.SetDefault("no-downloads", false)
	viper.SetDefault("formats", []string{".epub"})
	viper.SetDefault("encode", ".yaml")
}

func processMetadata(books []cdb.Book) {
	for _, b := range books {
		m := b.StringMap()
		if !ao3.DontSave() {
			name := casing.Snake(b.Title)
			err := writeMetaFile(m, name)
			if err != nil {
				log.Fatal(err)
			}
			if ao3.IsPodfic() || ffmeta {
				err := writeFFMeta(m, name)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		if !ao3.NoDownloads() {
			downloadFormats(b)
		}
		//err := b.Print(enc, true)
		//if err != nil {
		//log.Fatal(err)
		//}
	}
}

func writeFFMeta(r map[string]any, name string) error {
	ff := audbk.NewFFMeta()
	audbk.BookToFFMeta(ff, r)

	ffm, err := os.Create(name + ".ini")
	if err != nil {
		return fmt.Errorf("write init error: %w\n", err)
	}
	defer ffm.Close()
	ff.WriteTo(ffm)

	return nil
}

func writeMetaFile(r map[string]any, name string) error {
	var err error

	enc := ao3.Encode()

	if _, ok := r["formats"]; ok {
		delete(r, "formats")
	}

	mf, err := os.Create(name + enc)
	defer mf.Close()

	switch enc {
	case ".yaml":
		err = yaml.NewEncoder(mf).Encode(r)
	case ".json":
		err = json.NewEncoder(mf).Encode(r)
	case ".toml":
		err = toml.NewEncoder(mf).Encode(r)
	}

	if err != nil {
		return fmt.Errorf("write meta file err %w\n", err)
	}
	return nil
}

func downloadFormats(b cdb.Book) {
	for _, f := range b.Formats {
		for _, ext := range ao3.Formats() {
			if strings.Contains(f, ext) {
				fmt.Printf("downloading %s\n", b.Title+ext)
				name := casing.Snake(b.Title) + ext
				ao3.DownloadWork(f, name)
			}
		}
	}
}
