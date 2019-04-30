package do_state

import (
	"context"
	"errors"
	"github.com/digitalocean/godo"
	"github.com/harshpreet93/dopaas/do_auth"
	"log"
	"os"
)

type projectState struct {
	numDroplets int
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
		opt.Page = page + 1
	}

	return nil, errors.New("cannot find project with ID " + projectId)
}

func extractProjectResourceInfo(project *godo.Project) (*projectState, error) {
	return nil, nil
}

func GetState(projectId string) (*projectState, error) {
	log.Println("getting current state of project ", projectId)

	project, err := getProject(projectId)

	if err != nil {
		log.Println("error getting project ", err)
	}

	currState, err := extractProjectResourceInfo(project)

	if err != nil {
		log.Println("error getting current state", err)
		os.Exit(1)
	}
	return currState, nil
}
