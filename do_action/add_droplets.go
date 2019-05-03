package do_action

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/do_auth"
	"log"
	"os"
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
	droplets, _, err := do_auth.Auth().Droplets.CreateMultiple(ctx, dropletMultiCreateRequest)
	if err != nil {
		log.Println("Error adding droplets ", err)
	}
	assignDropletsToProject(droplets, conf.GetConfig().GetString("project_id"))
	return nil
}

func assignDropletsToProject(droplets []godo.Droplet, projectId string) {
		ctx := context.Background()
		projectsService := &godo.ProjectsServiceOp{do_auth.Auth() }
		_, _, err := projectsService.AssignResources(ctx, projectId, droplets)
		if err != nil {
			log.Println("error assigning resources to project ", err)
			os.Exit(1)
		}
}