package audit

// Observer defines a receiver of audit events
type Observer interface {
	// Send delivers an audit event
	Send(event Event) error
}
