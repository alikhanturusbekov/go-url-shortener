package audit

// Noop is a no-op audit event publisher
type Noop struct{}

// NewNoop creates a new Noop instance
func NewNoop() *Noop {
	return &Noop{}
}

// Notify performs no action.
func (n *Noop) Notify(event Event) {}

// Close performs no action.
func (n *Noop) Close() error {
	return nil
}
