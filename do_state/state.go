package do_state

type projectState struct {
	numDroplets int
}

func GetState(projectName string) (*projectState, error) {
	return &projectState{
		numDroplets: 3,
	}, nil
}
