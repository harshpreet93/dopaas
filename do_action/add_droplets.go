package do_action

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/harshpreet93/dopaas/do_auth"
	"log"
)

type AddDroplets struct {
	DesiredNum int
	SizeSlug   string
	Region     string
	ImageSlug  string
}

func (a AddDroplets) Execute() error {
	ctx := context.Background()
	dropletMultiCreateRequest := &godo.DropletMultiCreateRequest{Names: []string{"1", "2", "3"}, Region: a.Region, Size: a.SizeSlug, Image: godo.DropletCreateImage{Slug: a.ImageSlug}}
	_, _, err := do_auth.Auth().Droplets.CreateMultiple(ctx, dropletMultiCreateRequest)
	if err != nil {
		log.Println("Error adding droplets ", err)
	}
	return nil
}
