package do_action

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/do_auth"
	"log"
	"strconv"
)

type AddDroplets struct {
	DesiredNum int
	SizeSlug   string
	Region     string
	ImageSlug  string
}

func generateDropletNames(numDesired int) []string {
	var names []string
	for i := 0; i < numDesired; i++ {
		names = append(names, conf.GetConfig().GetString("project_name")+"--"+(strconv.Itoa(i)))
	}
	return names
}

func (a AddDroplets) Execute() error {
	ctx := context.Background()
	dropletMultiCreateRequest := &godo.DropletMultiCreateRequest{Names: generateDropletNames(conf.GetConfig().GetInt("NumDroplets")), Region: a.Region, Size: a.SizeSlug, Image: godo.DropletCreateImage{Slug: a.ImageSlug}}
	_, _, err := do_auth.Auth().Droplets.CreateMultiple(ctx, dropletMultiCreateRequest)
	if err != nil {
		log.Println("Error adding droplets ", err)
	}
	return nil
}
