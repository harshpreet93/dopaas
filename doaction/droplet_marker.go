package doaction

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/harshpreet93/dopaas/doauth"
	"time"
)

type DropletMarker struct {
	dropletID int
	filename string
	info string
}

func (d DropletMarker) Execute(runID string) error {

}
