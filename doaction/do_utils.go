package doaction

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/harshpreet93/dopaas/doauth"
	"time"
)

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
		return IP, nil

	}
	return "", err
}
