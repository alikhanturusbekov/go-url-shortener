package audit

import (
	"context"
	"log"
	"sync"

	"go.uber.org/zap"

	"github.com/alikhanturusbekov/go-url-shortener/pkg/logger"
)

// Service implements Publisher and dispatches events to registered observers
type Service struct {
	observers []Observer
	ch        chan Event
	wg        sync.WaitGroup
	closeOnce sync.Once
	ctx       context.Context
}

// NewService creates a new audit service with a buffered channel
func NewService(ctx context.Context, buffer int) *Service {
	s := &Service{
		ch:  make(chan Event, buffer),
		ctx: ctx,
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

	for {
		select {
		case <-s.ctx.Done():
			for event := range s.ch {
				s.dispatch(event)
			}
			return

		case event, ok := <-s.ch:
			if !ok {
				return
			}
			s.dispatch(event)
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

// dispatch notifies all the observers for the event
func (s *Service) dispatch(event Event) {
	for _, observer := range s.observers {
		if err := observer.Send(event); err != nil {
			logger.Log.Error("audit send error", zap.Error(err))
		}
	}
}
