package do_action

type Action interface {
	Execute() error
}
