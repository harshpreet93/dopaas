package conf

import (
	"bytes"
	"github.com/harshpreet93/dopaas/errorcheck"
	"github.com/spf13/viper"
	"io/ioutil"
)

type DesiredState struct {
	NumDroplets int
	SizeSlug    string
	Region      string
	ImageSlug   string
}

// marshal config file into a struct here
func GetConfig() *viper.Viper {
	conf := viper.New()
	conf.SetConfigType("yaml")
	dat, err := ioutil.ReadFile("app.yml")
	errorcheck.ExitOn(err, "error reading yaml config ")
	conf.ReadConfig(bytes.NewBuffer(dat))
	return conf
}

func GetDesiredState() (*DesiredState, error) {
	desiredState := &DesiredState{}
	desiredState.NumDroplets = GetConfig().GetInt("NumDroplets")
	desiredState.Region = GetConfig().GetString("Region")
	desiredState.ImageSlug = GetConfig().GetString("ImageSlug")
	desiredState.SizeSlug = GetConfig().GetString("SizeSlug")
	return desiredState, nil
}
