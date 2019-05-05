package do_action

import (
	"context"
	"github.com/fatih/color"
	"github.com/harshpreet93/dopaas/do_auth"
	"log"
	"os"
)

type DestroyDropletsAction struct {
	DropletID int
}

func (d DestroyDropletsAction) Execute(runID string) error {
	d.Print(false)
	response, err := do_auth.Auth().Droplets.Delete(context.Background(), d.DropletID)
	if err != nil {
		log.Println("error destroying droplet ", response, err)
		os.Exit(1)
	}
	return nil
}

func (d DestroyDropletsAction) Print(dryRun bool) {
	color.Red("--- Destroying droplet with ID %s", d.DropletID)
}
