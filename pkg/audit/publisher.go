package audit

type Publisher interface {
	Notify(event Event)
	Close()
}
