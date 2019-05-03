package do_state

import (
	"context"
	"errors"
	"github.com/digitalocean/godo"
	"github.com/harshpreet93/dopaas/do_auth"
	"log"
	"os"
	"strconv"
	"strings"
)

type ProjectState struct {
	Droplets []*godo.Droplet
}

func getProject(projectId string) (*godo.Project, error) {
	client := do_auth.Auth()
	ctx := context.Background()
	// create options. initially, these will be blank
	opt := &godo.ListOptions{}

	for {
		projects, resp, err := client.Projects.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		for _, project := range projects {
			if project.ID == projectId {
				return &project, nil
			}
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
		log.Println("page is ", page)
		opt.Page = page + 1
	}

	return nil, errors.New("cannot find project with ID " + projectId)
}

func extractProjectResourceInfo(project *godo.Project) ([]*godo.Droplet, error) {
	client := do_auth.Auth()
	ctx := context.Background()
	opt := &godo.ListOptions{}
	var droplets []*godo.Droplet
	for {
		projectResources, resp, err := client.Projects.ListResources(ctx, project.ID, opt)

		if err != nil {
			return nil, err
		}

		for _, projectResource := range projectResources {
			if strings.HasPrefix(projectResource.URN, "do:droplet:") {
				dropletId, _ := strconv.Atoi(strings.TrimPrefix(projectResource.URN, "do:droplet:"))
				droplet, _, err := client.Droplets.Get(ctx, dropletId)
				if err != nil {
					os.Exit(1)
				}
				droplets = append(droplets, droplet)
			}
		}
		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() || resp.Links.Pages.First == resp.Links.Pages.Last {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}
		// set the page we want for the next request
		opt.Page = page + 1
	}

	return droplets, nil
}

func GetState(projectId string) (*ProjectState, error) {
	project, err := getProject(projectId)
	if err != nil {
		os.Exit(1)
	}
	currState, err := extractProjectResourceInfo(project)
	if err != nil {
		os.Exit(1)
	}
	projectState := &ProjectState{droplets: currState}
	return projectState, nil
}
