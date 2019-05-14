package doaction

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/fatih/color"
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/doauth"
	"github.com/harshpreet93/dopaas/errorcheck"
	"github.com/harshpreet93/dopaas/keyutil"
	"github.com/harshpreet93/dopaas/userdata"
	"strconv"
)

type AddDroplets struct {
	DesiredNum int
	SizeSlug   string
	Region     string
	ImageSlug  string
}

func generateDropletNames(numDesired int, runID string) []string {
	var names []string
	for i := 0; i < numDesired; i++ {
		names = append(names, conf.GetConfig().GetString("project_name")+"--"+(strconv.Itoa(i))+"--"+runID)
	}
	return names
}

func (a AddDroplets) Execute(runID string) error {
	a.Print(false)
	ctx := context.Background()
	dropletMultiCreateRequest := &godo.DropletMultiCreateRequest{
		Names:  generateDropletNames(a.DesiredNum, runID),
		Region: a.Region,
		Size:   a.SizeSlug,
		Image: godo.DropletCreateImage{
			Slug: a.ImageSlug,
		},
		Monitoring: true,
		SSHKeys:    []godo.DropletCreateSSHKey{{Fingerprint: keyutil.GetPubKeySignature()}},
		UserData:   userdata.Generate(),
	}
	droplets, _, err := doauth.Auth().Droplets.CreateMultiple(ctx, dropletMultiCreateRequest)
	if err != nil {
		color.Red("Error adding droplets ", err)
	}
	for _, droplet := range droplets {
		err = Transport{
			ID:           droplet.ID,
			ArtifactFile: conf.GetConfig().GetString("artifact_file"),
		}.Execute(runID)

		if err != nil {
			return err
		}

		err = Starter{
			ID: droplet.ID,
		}.Execute(runID)

		if err != nil {
			return err
		}

		err = Tagger{
			DropletId: droplet.ID,
			Tag:       "artifact_rev_" + GetFileSha(conf.GetConfig().GetString("artifact_file")),
		}.Execute(runID)
		errorcheck.ExitOn(err, "error tagging droplet")

		err = DropletMarker{
			dropletID: droplet.ID,
			Info:      GetFileSha(conf.GetConfig().GetString("artifact_file")),
			Filename:  "/root/artifact_sha",
		}.Execute(runID)
		errorcheck.ExitOn(err, "error marking droplet with sha")
	}
	return nil
}

func (a AddDroplets) Print(dryRun bool) {
	prefix := "Creating"
	if dryRun {
		prefix = "Would create"
	}
	color.Green("+++ "+prefix+" %d droplets in %s with size %s and image %s",
		a.DesiredNum, a.Region, a.SizeSlug, a.ImageSlug)
}
