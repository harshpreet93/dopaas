package conf

import (
	"bytes"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
)

type DesiredState struct {
	NumDroplets int
	SizeSlug    string
	Region      string
	ImageSlug   string
}

// marshal config file into a struct here
func GetConfig() *viper.Viper {
	os.Getwd()
	conf := viper.New()
	conf.SetConfigType("yaml")
	dat, err := ioutil.ReadFile("app.yml")
	if err != nil {
		log.Println("error reading yaml config ", err)
		os.Exit(1)
	}
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
