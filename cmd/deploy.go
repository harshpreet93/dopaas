package cmd

import (
	"fmt"
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/do_action"
	"github.com/harshpreet93/dopaas/do_state"
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
	projectId := conf.GetConfig().Get("project_name").(string)
	currState, err := do_state.GetState(projectId)
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
	for _, action := range actions {
		log.Println("action required ", *action)
		(*action).Execute()
	}
}

func diff(state *do_state.ProjectState, desiredState *conf.DesiredState) ([]*do_action.Action, error) {
	var actions []*do_action.Action

	if desiredState.NumDroplets < len(state.Droplets) {
		// TODO: if something must be destroyed and not all can be destroyed......which should be preferred for destruction?
		var destroy do_action.Action = &do_action.DestroyDropletsAction{}
		actions = append(actions, &destroy)
	}

	if desiredState.NumDroplets > len(state.Droplets) {
		var add do_action.Action = &do_action.AddDroplets{DesiredNum: desiredState.NumDroplets - len(state.Droplets),
			ImageSlug: conf.GetConfig().GetString("ImageSlug"),
			SizeSlug:  conf.GetConfig().GetString("SizeSlug"),
			Region:    conf.GetConfig().GetString("Region")}
		actions = append(actions, &add)
	}
	return actions, nil
}
