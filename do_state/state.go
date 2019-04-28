package do_state

import (
	"context"
	"github.com/harshpreet93/dopaas/do_auth"
	"log"
)

type projectState struct {
	numDroplets int
}

func GetState(projectName string) (*projectState, error) {
	log.Println("getting current state of project ", projectName)
	client := do_auth.Auth()
	ctx := context.Background()
	_, _, err := client.Projects.Get(ctx, projectName)

	if err != nil {
		log.Println("error getting project ", projectName, err)
	}

	return &projectState{
		numDroplets: 3,
	}, nil
}
