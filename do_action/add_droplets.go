package do_action

import (
	"bytes"
	"context"
	"github.com/digitalocean/godo"
	"github.com/fatih/color"
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/do_auth"
	"github.com/harshpreet93/dopaas/error_check"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"text/template"
)

type AddDroplets struct {
	DesiredNum int
	SizeSlug   string
	Region     string
	ImageSlug  string
}

func generateDropletNames(numDesired int, runID string) []string {
	var names []string
	for i := 0; i < numDesired; i++ {
		names = append(names, conf.GetConfig().GetString("project_name")+"--"+(strconv.Itoa(i))+"--"+runID)
	}
	return names
}

func (a AddDroplets) Execute(runID string) error {
	a.Print(false)
	ctx := context.Background()
	sshPubKeyFile, err := homedir.Expand("~/.ssh/id_rsa.pub")
	error_check.ExitOn(err, "error getting pub key in ~/.ssh/id_rsa.pub")
	sshPubKeyContents, err := ioutil.ReadFile(sshPubKeyFile)
	error_check.ExitOn(err, "Error getting pub key file contents")
	log.Println("pubKeyContents ", string( sshPubKeyContents ))

	userData := `#!/bin/sh
				mkdir -p ~/.ssh && touch ~/.ssh/authorized_keys
				chmod 700 ~/.ssh && chmod 600 ~/.ssh/authorized_keys
				echo {{.pubKey}} >> ~/.ssh/authorized_keys
				# add key
				# create ssh user
				# create app user`

	tmpl, err := template.New("userdata").Parse(userData)
	error_check.ExitOn(err, "error creating userdata template")
	tmplVars := template.FuncMap{
		"pubKey": strings.TrimSpace(string(sshPubKeyContents)),
	}
	compiledUserData := &bytes.Buffer{}
	err = tmpl.Execute(compiledUserData,tmplVars)
	error_check.ExitOn(err, "error compiling userdata template")
	log.Println("compiled user data is ", compiledUserData.String())
	dropletMultiCreateRequest := &godo.DropletMultiCreateRequest{
		Names:  generateDropletNames(a.DesiredNum, runID),
		Region: a.Region,
		Size:   a.SizeSlug,
		Image: godo.DropletCreateImage{
			Slug: a.ImageSlug,
		},
		Monitoring: true,
		UserData:   compiledUserData.String(),
	}
	_, _, err = do_auth.Auth().Droplets.CreateMultiple(ctx, dropletMultiCreateRequest)
	if err != nil {
		color.Red("Error adding droplets ", err)
	}
	return nil
}

func (a AddDroplets) Print(dryRun bool) {
	prefix := "Creating"
	if dryRun {
		prefix = "Would create"
	}
	color.Green("+++ "+prefix+" %d droplets in %s with size %s and image %s",
		a.DesiredNum, a.Region, a.SizeSlug, a.ImageSlug)
}
