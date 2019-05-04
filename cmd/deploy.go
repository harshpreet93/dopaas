package cmd

import (
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/do_action"
	"github.com/harshpreet93/dopaas/do_state"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an app to digitalocean",
	Run:   do,
}

func do(cmd *cobra.Command, args []string) {
	projectName := conf.GetConfig().Get("project_name").(string)
	currState, err := do_state.GetState(projectName)
	if err != nil {
		log.Println("Error getting current state. exiting", err)
		os.Exit(1)
	}
	desiredState, err := conf.GetDesiredState()
	fmt.Println(desiredState, currState)
	actions, err := diff(currState, desiredState)
	if err != nil {
		log.Println("Error calculating actions to get to desired state ", err)
		os.Exit(1)
	}
	runID := uuid.NewV4()

	for _, action := range actions {
		log.Println("action required ", *action)
		(*action).Execute(runID.String())
	}
}

func diff(state *do_state.ProjectState, desiredState *conf.DesiredState) ([]*do_action.Action, error) {
	var actions []*do_action.Action
	numDropletsToBeCreated := desiredState.NumDroplets
	var toBeDestroyed []*godo.Droplet
	log.Println("num Existing droplets", len(state.Droplets))
	for _, droplet := range state.Droplets {
		log.Println("found existing droplet ", droplet.SizeSlug == desiredState.SizeSlug, droplet.Region.Slug == desiredState.Region, droplet.Image.Slug == desiredState.ImageSlug)
		if droplet.SizeSlug == desiredState.SizeSlug &&
			droplet.Region.Slug == desiredState.Region &&
			droplet.Image.Slug == desiredState.ImageSlug {
			numDropletsToBeCreated--
		} else {
			toBeDestroyed = append(toBeDestroyed, droplet)
		}
	}
	if numDropletsToBeCreated > 0 {
		var add do_action.Action = &do_action.AddDroplets{
			DesiredNum: numDropletsToBeCreated,
			ImageSlug:  conf.GetConfig().GetString("ImageSlug"),
			SizeSlug:   conf.GetConfig().GetString("SizeSlug"),
			Region:     conf.GetConfig().GetString("Region"),
		}
		actions = append(actions, &add)
	}

	return actions, nil
}
