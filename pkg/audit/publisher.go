package audit

// Publisher defines an audit event publisher
type Publisher interface {
	// Notify publishes an audit event
	Notify(event Event)

	// Close closes publisher
	Close() error
}
