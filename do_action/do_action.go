package do_action

type Action interface {
	Execute(runID string) error
}
