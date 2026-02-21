package audit

import (
	"log"
	"sync"
)

// Service implements Publisher and dispatches events to registered observers
type Service struct {
	observers []Observer
	ch        chan Event
	wg        sync.WaitGroup
	closeOnce sync.Once
}

// NewService creates a new audit service with a buffered channel
func NewService(buffer int) *Service {
	s := &Service{
		ch: make(chan Event, buffer),
	}

	s.wg.Add(1)
	go s.worker()

	return s
}

// Register adds an observer to receive audit events
func (s *Service) Register(o Observer) {
	s.observers = append(s.observers, o)
}

// Notify enqueues an audit event for asynchronous delivery
func (s *Service) Notify(event Event) {
	select {
	case s.ch <- event:
	default:
		log.Println("audit buffer full, dropping event")
	}
}

// worker processes events and forwards them to observers
func (s *Service) worker() {
	defer s.wg.Done()

	for event := range s.ch {
		for _, obs := range s.observers {
			if err := obs.Send(event); err != nil {
				log.Printf("audit send error: %v", err)
			}
		}
	}
}

// Close stops the service and waits for pending events to be processed
func (s *Service) Close() error {
	s.closeOnce.Do(func() {
		close(s.ch)
	})
	s.wg.Wait()
	return nil
}
