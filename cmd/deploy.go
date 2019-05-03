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
	projectId := conf.GetConfig().Get("project_id").(string)
	currState, err := do_state.GetState(projectId)
	if err != nil {
		log.Println("Error getting current state. exiting", err)
		os.Exit(1)
	}
	desiredState, err := conf.GetDesiredState()
	fmt.Println(desiredState, currState)
	diff(currState, desiredState)
}

func diff(state *do_state.ProjectState, desiredState *conf.DesiredState) ([]*do_action.Action, error) {
	var actions []*do_action.Action
	if desiredState.NumDroplets > len( state.Droplets ) {
		actions := append( actions, &do_action.CreateDropletsAction{} )
	}
	return nil, nil
}