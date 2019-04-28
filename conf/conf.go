package conf

import "os"

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
func getConfig() {
	os.Getwd()
}
