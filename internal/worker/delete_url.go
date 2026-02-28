package worker

import (
	"context"
	"time"

	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
	"github.com/alikhanturusbekov/go-url-shortener/internal/repository"
	"github.com/alikhanturusbekov/go-url-shortener/pkg/logger"
)

// DeleteURLWorker processes URL deletion tasks asynchronously
type DeleteURLWorker struct {
	repository repository.URLRepository
	in         chan model.DeleteURLTask
}

// NewDeleteURLWorker creates a new DeleteURLWorker instance
func NewDeleteURLWorker(
	repository repository.URLRepository,
	bufferSize int,
) *DeleteURLWorker {
	return &DeleteURLWorker{
		repository: repository,
		in:         make(chan model.DeleteURLTask, bufferSize),
	}
}

// Enqueue adds a deletion task to the worker queue
func (w *DeleteURLWorker) Enqueue(task model.DeleteURLTask) {
	w.in <- task
}

// Run starts the worker loop and processes tasks in batches
func (w *DeleteURLWorker) Run(ctx context.Context) {
	const (
		maxBatchSize = 100
		flushTimeout = 500 * time.Millisecond
	)

	ticker := time.NewTicker(flushTimeout)
	defer ticker.Stop()

	buffer := make([]model.DeleteURLTask, 0, maxBatchSize)

	flush := func() {
		if len(buffer) == 0 {
			return
		}

		grouped := make(map[string][]string)
		for _, task := range buffer {
			grouped[task.UserID] = append(grouped[task.UserID], task.Short)
		}

		for userID, urls := range grouped {
			if err := w.repository.DeleteByShorts(ctx, userID, urls); err != nil {
				logger.Log.Error("could not delete user URLs:" + err.Error())
			}
		}

		buffer = buffer[:0]
	}

	for {
		select {
		case <-ctx.Done():
			flush()
			return

		case task := <-w.in:
			buffer = append(buffer, task)

			if len(buffer) >= maxBatchSize {
				flush()
			}

		case <-ticker.C:
			flush()
		}
	}
}
