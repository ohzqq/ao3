package ao3

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("podfic", false)
	viper.SetDefault("no-save", false)
	viper.SetDefault("no-downloads", false)
	viper.SetDefault("formats", []string{".epub"})
	viper.SetDefault("encode", []string{".yaml"})
}

func TestOpts(t *testing.T) {
	fmt.Printf(optFmt,
		IsPodfic(),
		NoDownloads(),
		DontSave(),
		Formats(),
		Encode(),
	)
}

const optFmt = `
IsPodfic: %v
NoDownloads: %v
DontSave: %v
Formats: %v
Encode: %v
`
