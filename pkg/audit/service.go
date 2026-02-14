package audit

import (
	"log"
	"sync"
)

type Service struct {
	observers []Observer
	ch        chan Event
	wg        sync.WaitGroup
}

func NewService(buffer int) *Service {
	s := &Service{
		ch: make(chan Event, buffer),
	}

	s.wg.Add(1)
	go s.worker()

	return s
}

func (s *Service) Register(o Observer) {
	s.observers = append(s.observers, o)
}

func (s *Service) Notify(event Event) {
	select {
	case s.ch <- event:
	default:
		log.Println("audit buffer full, dropping event")
	}
}

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

func (s *Service) Close() {
	close(s.ch)
	s.wg.Wait()
}
