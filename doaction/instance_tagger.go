package doaction

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/harshpreet93/dopaas/doauth"
)

type Tagger struct {
	DropletId int
	Tag string
}

func (a Tagger) Execute(runID string) error {
	client := doauth.Auth()
	ctx := context.TODO()
	tagRequest := &godo.TagResourcesRequest{
		Resources: []godo.Resource{
			{ID: string(a.DropletId), Type: godo.DropletResourceType},
		},
	}
	_, _, err := client.Tags.Create(ctx, &godo.TagCreateRequest{Name: a.Tag})
	if err != nil {
		return err
	}
	_, err = client.Tags.TagResources(ctx, a.Tag, tagRequest)
	return err
}