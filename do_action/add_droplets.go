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

func generateDropletNames(numDesired int, runID string) []string {
	var names []string
	for i := 0; i < numDesired; i++ {
		names = append(names, conf.GetConfig().GetString("project_name")+"--"+(strconv.Itoa(i))+"--"+runID)
	}
	return names
}

func (a AddDroplets) Execute(runID string) error {
	ctx := context.Background()
	dropletMultiCreateRequest := &godo.DropletMultiCreateRequest{
		Names:  generateDropletNames(a.DesiredNum, runID),
		Region: a.Region,
		Size:   a.SizeSlug,
		Image: godo.DropletCreateImage{
			Slug: a.ImageSlug,
		},
	}
	_, _, err := do_auth.Auth().Droplets.CreateMultiple(ctx, dropletMultiCreateRequest)
	if err != nil {
		log.Println("Error adding droplets ", err)
	}
	return nil
}
