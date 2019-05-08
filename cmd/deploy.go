package cmd

import (
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/doaction"
	"github.com/harshpreet93/dopaas/dostate"
	"github.com/harshpreet93/dopaas/errorcheck"
	"github.com/harshpreet93/dopaas/predeploy"
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
	currState, err := dostate.GetState(projectName)
	errorcheck.ExitOn(err, "Error getting current state. exiting")
	predeploy.Execute()
	desiredState, err := conf.GetDesiredState()
	actions, err := diff(currState, desiredState)
	errorcheck.ExitOn(err, "Error calculating actions to get to desired state ")
	runID := uuid.NewV4()
	for _, action := range actions {
		if !dryRun {
			(*action).Execute(runID.String())
		} else {
			(*action).Print(dryRun)
		}
	}
}

func diff(state *dostate.ProjectState, desiredState *conf.DesiredState) ([]*doaction.Action, error) {
	var actions []*doaction.Action
	numDropletsToBeCreated := desiredState.NumDroplets
	for _, droplet := range state.Droplets {
		if droplet.SizeSlug == desiredState.SizeSlug &&
			droplet.Region.Slug == desiredState.Region &&
			droplet.Image.Slug == desiredState.ImageSlug {
			numDropletsToBeCreated--
		} else {
			var destroy doaction.Action = &doaction.DestroyDropletsAction{DropletID: droplet.ID}
			actions = append(actions, &destroy)
		}

		if numDropletsToBeCreated < 0 {
			var destroy doaction.Action = &doaction.DestroyDropletsAction{DropletID: droplet.ID}
			actions = append(actions, &destroy)
		}
	}
	if numDropletsToBeCreated > 0 {
		var add doaction.Action = &doaction.AddDroplets{
			DesiredNum: numDropletsToBeCreated,
			ImageSlug:  conf.GetConfig().GetString("ImageSlug"),
			SizeSlug:   conf.GetConfig().GetString("SizeSlug"),
			Region:     conf.GetConfig().GetString("Region"),
		}
		actions = append(actions, &add)
	}

	return actions, nil
}
