package do_action

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/fatih/color"
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/do_auth"
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
	a.Print(false)
	ctx := context.Background()

	userData := `#!/bin/sh
				# add key
				# create ssh user
				# create app user`

	dropletMultiCreateRequest := &godo.DropletMultiCreateRequest{
		Names:  generateDropletNames(a.DesiredNum, runID),
		Region: a.Region,
		Size:   a.SizeSlug,
		Image: godo.DropletCreateImage{
			Slug: a.ImageSlug,
		},
		Monitoring: true,
		UserData: userData,
	}
	_, _, err := do_auth.Auth().Droplets.CreateMultiple(ctx, dropletMultiCreateRequest)
	if err != nil {
		color.Red("Error adding droplets ", err)
	}
	return nil
}

func (a AddDroplets) Print(dryRun bool) {
	prefix := "Creating"
	if dryRun {
		prefix = "Would create"
	}
	color.Green("+++ "+prefix+" %d droplets in %s with size %s and image %s",
		a.DesiredNum, a.Region, a.SizeSlug, a.ImageSlug)
}
