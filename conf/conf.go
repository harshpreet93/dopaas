package conf

import (
	"bytes"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
)

//look for app.yml or app.yaml files in the run dir and load if available

//test: []
//build: ["", "", ""]
//
//#deploy info
//domain: syneo.io
//lb: true
//droplet_type: 'blah'
//num_droplets: 3
//dc: 'blah'
//key_location: ''
//start: [""]

type config struct {
	test         string
	build        string
	domain       string
	lb           bool
	droplet_type string
	droplet_num  int64
	dc           string
	start        string
}

// marshal config file into a struct here
func GetConfig() *viper.Viper {
	os.Getwd()
	conf := viper.New()
	conf.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	dat, err := ioutil.ReadFile("app.yml")
	if err != nil {
		log.Println("error reading yaml config ", err)
		os.Exit(1)
	}
	conf.ReadConfig(bytes.NewBuffer(dat))
	log.Println(dat)
	log.Println("project ID is ", viper.Get("project_id"))
	return conf
}
