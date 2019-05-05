package cmd

import (
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/do_action"
	"github.com/harshpreet93/dopaas/do_state"
	"github.com/harshpreet93/dopaas/error_check"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

var dryRun bool

func init() {
	deployCmd.Flags().BoolVar(&dryRun, "dryrun", false, "Show what deploy would do instead of actually doing it")
	rootCmd.AddCommand(deployCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an app to digitalocean",
	Run:   do,
}

func do(cmd *cobra.Command, args []string) {
	projectName := conf.GetConfig().Get("project_name").(string)
	currState, err := do_state.GetState(projectName)
	error_check.ExitOn(err, "Error getting current state. exiting")
	desiredState, err := conf.GetDesiredState()
	actions, err := diff(currState, desiredState)
	error_check.ExitOn(err, "Error calculating actions to get to desired state ")
	runID := uuid.NewV4()
	for _, action := range actions {
		if !dryRun {
			(*action).Execute(runID.String())
		} else {
			(*action).Print(dryRun)
		}
	}
}

func diff(state *do_state.ProjectState, desiredState *conf.DesiredState) ([]*do_action.Action, error) {
	//TODO: Fix a bug, that can be reproduced by trying to downsize from, for instance, 3 to 1 droplet......only one droplet is destroyed..
	var actions []*do_action.Action
	numDropletsToBeCreated := desiredState.NumDroplets
	for _, droplet := range state.Droplets {

		if droplet.SizeSlug == desiredState.SizeSlug &&
			droplet.Region.Slug == desiredState.Region &&
			droplet.Image.Slug == desiredState.ImageSlug {
			numDropletsToBeCreated--
		} else {
			var destroy do_action.Action = &do_action.DestroyDropletsAction{DropletID: droplet.ID}
			actions = append(actions, &destroy)
		}

		if numDropletsToBeCreated < 0 {
			var destroy do_action.Action = &do_action.DestroyDropletsAction{DropletID: droplet.ID}
			actions = append(actions, &destroy)
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
