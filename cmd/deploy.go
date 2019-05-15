package cmd

import (
	"fmt"
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/doaction"
	"github.com/harshpreet93/dopaas/dostate"
	"github.com/harshpreet93/dopaas/errorcheck"
	"github.com/harshpreet93/dopaas/predeploy"
	"github.com/kyokomi/emoji"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"log"
	"strings"
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
	success_message := emoji.Sprint(":rocket: :rocket: :rocket: SUCCESS! :rocket: :rocket: :rocket:")
	fmt.Println(success_message)

}

func diff(state *dostate.ProjectState, desiredState *conf.DesiredState) ([]*doaction.Action, error) {
	var actions []*doaction.Action
	numDropletsToBeCreated := desiredState.NumDroplets
	var dropletsToBeDestroyed []int
	for _, droplet := range state.Droplets {
		if droplet.SizeSlug == desiredState.SizeSlug &&
			droplet.Region.Slug == desiredState.Region &&
			droplet.Image.Slug == desiredState.ImageSlug {
			numDropletsToBeCreated--
		} else {
			var destroy doaction.Action = &doaction.DestroyDropletsAction{DropletID: droplet.ID}
			actions = append(actions, &destroy)
			dropletsToBeDestroyed = append(dropletsToBeDestroyed, droplet.ID)
		}

		if numDropletsToBeCreated < 0 {
			var destroy doaction.Action = &doaction.DestroyDropletsAction{DropletID: droplet.ID}
			actions = append(actions, &destroy)
			dropletsToBeDestroyed = append(dropletsToBeDestroyed, droplet.ID)
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
	for _, droplet := range state.Droplets {
		if !contains(dropletsToBeDestroyed, droplet.ID) &&
			strings.TrimSpace(doaction.GetDropletArtifactSha(droplet.ID)) !=
				strings.TrimSpace(doaction.GetFileSha(conf.GetConfig().GetString("artifact_file"))) {
			var transport doaction.Action = doaction.Action(&doaction.Transport{
				ID:           droplet.ID,
				ArtifactFile: conf.GetConfig().GetString("artifact_file"),
			})
			var starter doaction.Action = doaction.Starter{
				ID: droplet.ID,
			}
			var marker doaction.Action = doaction.DropletMarker{
				DropletID: droplet.ID,
				Info:      doaction.GetFileSha(conf.GetConfig().GetString("artifact_file")),
				Filename:  "/root/artifact_sha",
			}
			actions = append(actions, &transport, &starter, &marker)
			log.Println("droplet needs new version of artifact ", droplet.ID, doaction.GetDropletArtifactSha(droplet.ID), doaction.GetFileSha(conf.GetConfig().GetString("artifact_file")))
		}
		log.Printf("current sha on %s is %s", string(droplet.ID), doaction.GetDropletArtifactSha(droplet.ID))

	}
	return actions, nil
}

func contains(a []int, x int) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
