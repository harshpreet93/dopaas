package do_state

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/harshpreet93/dopaas/do_auth"
	"github.com/harshpreet93/dopaas/error_check"
	"strings"
)

type ProjectState struct {
	Droplets []*godo.Droplet
}

func GetState(projectName string) (*ProjectState, error) {
	currState := getDropletsForProject(projectName)
	projectState := &ProjectState{Droplets: currState}
	return projectState, nil
}

func getAllDroplets(ctx context.Context, client *godo.Client) ([]godo.Droplet, error) {
	// create a list to hold our droplets
	var list []godo.Droplet

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := client.Droplets.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		// append the current page's droplets to our list
		for _, d := range droplets {
			list = append(list, d)
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return list, nil
}

func getDropletsForProject(projectName string) []*godo.Droplet {
	var currState []*godo.Droplet
	allDroplets, err := getAllDroplets(context.Background(), do_auth.Auth())
	error_check.ExitOn(err, "error getting current project state")
	for _, droplet := range allDroplets {
		if strings.HasPrefix(droplet.Name, projectName) {
			currState = append(currState, &droplet)
		}
	}
	return currState
}
