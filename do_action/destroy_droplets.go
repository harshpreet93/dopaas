package do_action

import "log"

type DestroyDropletsAction struct {
	DropletID string
}

func (d DestroyDropletsAction) Execute(runID string) error {
	//TODO
	log.Println("destroying ", d.DropletID)
	return nil
}
