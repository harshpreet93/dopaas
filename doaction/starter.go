package doaction

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/doauth"
	"github.com/harshpreet93/dopaas/errorcheck"
	"github.com/sfreiberg/simplessh"
	"log"
	"time"
)

type Starter struct {
	ID int
}

func (a Starter) Execute(runID string) error {
	//ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//defer cancel()
	return a.executeWithTimeout( runID)
}

func tryToGetIPForId(ID int) (string, error) {
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
		log.Printf("found IP %s for ID %d", IP, ID)
		return IP, nil

	}
	return "", err
}

func (a Starter) executeWithTimeout(runID string) error {

	d := time.Now().Add(30 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)

	defer cancel()

	select {
		case <-ctx.Done(): cancel()
	}

	ip, err := tryToGetIPForId(a.ID)
	errorcheck.ExitOn(err, "Error getting IP for droplet id")
	client, err := simplessh.ConnectWithKeyFile(ip+":22", "root", "")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	output, err := client.Exec("cd /root && " + conf.GetConfig().GetString("start"))
	log.Println("start script output", output)
	client.Close()
	return err
}
