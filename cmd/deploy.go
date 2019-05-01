package cmd

import (
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/do_state"
	"github.com/spf13/cobra"
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
	project_id := conf.GetConfig().Get("project_id").(string)
	do_state.GetState(project_id)
}
