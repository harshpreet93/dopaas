package doaction

import (
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/errorcheck"
	"github.com/sfreiberg/simplessh"
	"time"
)

type Starter struct {
	ID int
}

func (a Starter) Print(dryRun bool) {

}

func (a Starter) Execute(runID string) error {
	done := make(chan error)
	go a.executeWithTimeout(runID, done)
	select {
	case err := <-done:
		return err
	case <-time.After(30 * time.Second):
	}
	close(done)
	return nil
}

func (a Starter) executeWithTimeout(runID string, done chan error) {
	ip, err := tryToGetIPForId(a.ID)
	errorcheck.ExitOn(err, "Error getting IP for droplet id")
	client, err := simplessh.ConnectWithKeyFile(ip+":22", "root", "")
	errorcheck.ExitOn(err, "error establishing connection to "+ip)
	defer client.Close()
	_, _ = client.Exec("kill -9 $(lsof -t -i:8080)")
	_, err = client.Exec("cd /root && " + conf.GetConfig().GetString("start"))
	done <- err
	close(done)
}
