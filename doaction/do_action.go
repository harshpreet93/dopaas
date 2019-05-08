package doaction

type Action interface {
	Execute(runID string) error
	Print(dryRun bool)
}
