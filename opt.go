package ao3

import "github.com/spf13/viper"

func IsPodfic() bool {
	return viper.GetBool("podfic")
}

func NoDownloads() bool {
	return viper.GetBool("no-downloads")
}

func DontSave() bool {
	return viper.GetBool("no-save")
}

func Formats() []string {
	return viper.GetStringSlice("formats")
}

func Encode() []string {
	return viper.GetStringSlice("encode")
}
