package doaction

import (
	"github.com/fatih/color"
	"github.com/harshpreet93/dopaas/errorcheck"
	"github.com/sfreiberg/simplessh"
	"strconv"
	"time"
)

type DropletMarker struct {
	DropletID int
	Filename  string
	Info      string
}

func (d DropletMarker) Print(dryRun bool) {
	prefix := "marking "
	if dryRun {
		color.Green("would mark ")
	}
	color.Green("++++++ "+prefix+strconv.Itoa(d.DropletID)+" with artifact sha")
}

func (d DropletMarker) Execute(runID string) error {
	d.Print(false)
	done := make(chan error)
	go d.executeWithTimeout(runID, done)
	select {
	case err := <-done:
		return err
	case <-time.After(30 * time.Second):
	}
	close(done)
	return nil

}

func (d DropletMarker) executeWithTimeout(runID string, done chan error) {
	ip, err := tryToGetIPForId(d.DropletID)
	errorcheck.ExitOn(err, "Error getting IP for droplet id")
	client, err := simplessh.ConnectWithKeyFile(ip+":22", "root", "")
	errorcheck.ExitOn(err, "error establishing connection to "+ip)
	defer client.Close()
	_, err = client.Exec("set -eo pipefail; echo " + d.Info + " > " + d.Filename)
	done <- err
	close(done)
}
