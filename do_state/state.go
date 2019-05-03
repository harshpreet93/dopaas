package do_state

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/harshpreet93/dopaas/do_auth"
	"log"
)

type ProjectState struct {
	Droplets []*godo.Droplet
}


func GetState(projectName string) (*ProjectState, error) {
	currState, err := getDropletsForProject(projectName)
	if err != nil {
		log.Println("error getting current state of project ", err)
	}
	projectState := &ProjectState{Droplets: currState}
	return projectState, nil
}

func getDropletsForProject(projectName string) ([]*godo.Droplet, error) {
	var currState []*godo.Droplet

	do_auth.Auth().Droplets.ListByTag(context.Background(), projectName)
	return currState, nil
}