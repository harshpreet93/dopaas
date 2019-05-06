package do_action

import (
	"context"
	"github.com/fatih/color"
	"github.com/harshpreet93/dopaas/do_auth"
	"github.com/harshpreet93/dopaas/error_check"
)

type DestroyDropletsAction struct {
	DropletID int
}

func (d DestroyDropletsAction) Execute(runID string) error {
	d.Print(false)
	_, err := do_auth.Auth().Droplets.Delete(context.Background(), d.DropletID)
	error_check.ExitOn(err, "error destroying droplet ")
	return nil
}

func (d DestroyDropletsAction) Print(dryRun bool) {
	prefix := "Destroying"
	if dryRun {
		prefix = "Would destroy"
	}
	color.Red("--- "+prefix+" droplet with ID %s", d.DropletID)
}
