package do_action

import (
	"context"
	"github.com/harshpreet93/dopaas/do_auth"
)

type KeyCreator struct {

}

func (k KeyCreator) Execute(runID string) {

	client := do_auth.Auth()
	ctx := context.Background()
	for

}

func (KeyCreator) Print(dryRun bool) {

}