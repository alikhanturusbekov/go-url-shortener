package audit

type Noop struct{}

func NewNoop() *Noop {
	return &Noop{}
}

func (n *Noop) Notify(event Event) {}

func (n *Noop) Close() {}
