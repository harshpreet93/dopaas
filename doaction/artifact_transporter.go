package doaction

import (
	"bytes"
	"context"
	"github.com/digitalocean/godo"
	"github.com/fatih/color"
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/doauth"
	"log"
	"os/exec"
	"strconv"
	"time"
)

type Transport struct {
	ID           int
	ArtifactFile string
}

func (t Transport) Execute(runID string) error {
	t.Print(false)
	return t.try()
}

func (t Transport) try() error {
	var err error = nil
	for i := 0; i < 5; i++ {
		IP, err := t.tryToGetIPForId(t.ID)
		if err != nil {
			time.Sleep(10000)
			continue
		}
		cmd := exec.Command("rsync", "-e", "ssh -o StrictHostKeyChecking=no", "-r", "--delete",
			conf.GetConfig().GetString("artifact_file"), "root@"+IP+":/root", "--timeout=60")
		var errOut bytes.Buffer
		cmd.Stderr = &errOut
		var stdOut bytes.Buffer
		cmd.Stdout = &stdOut
		err = cmd.Start()
		if err != nil {
			log.Print(err)
		}
		err = cmd.Wait()
		if err != nil {
			log.Println(errOut.String(), stdOut.String())
		}
		time.Sleep(3000)
	}
	return err
}

func (t Transport) tryToGetIPForId(ID int) (string, error) {
	ctx := context.Background()
	var err error
	var droplet *godo.Droplet
	for i := 0; i < 5; i++ {
		droplet, _, err = doauth.Auth().Droplets.Get(ctx, ID)
		if err != nil {
			time.Sleep(10000)
			continue
		}
		IP, err := droplet.PublicIPv4()
		if err != nil || IP == "" {
			time.Sleep(10000)
			continue
		}
		return IP, nil

	}
	return "", err
}

func (t Transport) Print(dryRun bool) {
	prefix := "transporting"
	if dryRun {
		color.Green("would transport")
	}
	color.Green("++++++ "+prefix+" artifact to "+strconv.Itoa(t.ID))
}
