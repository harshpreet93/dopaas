package doaction

import (
	"github.com/harshpreet93/dopaas/errorcheck"
	"github.com/sfreiberg/simplessh"
	"log"
	"time"
)

type DropletMarker struct {
	DropletID int
	Filename  string
	Info      string
}

func (d DropletMarker) Print(dryRun bool) {

}

func (d DropletMarker) Execute(runID string) error {
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
	output, err := client.Exec("echo " + d.Info + " > " + d.Filename)
	log.Println("marker script output", output)
	done <- err
	close(done)
}
