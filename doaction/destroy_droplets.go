package doaction

import (
	"context"
	"github.com/fatih/color"
	"github.com/harshpreet93/dopaas/doauth"
	"github.com/harshpreet93/dopaas/errorcheck"
)

type DestroyDropletsAction struct {
	DropletID int
}

func (d DestroyDropletsAction) Execute(runID string) error {
	d.Print(false)
	_, err := doauth.Auth().Droplets.Delete(context.Background(), d.DropletID)
	errorcheck.ExitOn(err, "error destroying droplet ")
	return nil
}

func (d DestroyDropletsAction) Print(dryRun bool) {
	prefix := "Destroying"
	if dryRun {
		prefix = "Would destroy"
	}
	color.Red("--- "+prefix+" droplet with ID %s", d.DropletID)
}
