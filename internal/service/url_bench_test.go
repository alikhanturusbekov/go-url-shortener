package service

import (
	"testing"

	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
	"github.com/alikhanturusbekov/go-url-shortener/internal/repository"
	"github.com/alikhanturusbekov/go-url-shortener/internal/worker"
	"github.com/alikhanturusbekov/go-url-shortener/pkg/audit"
)

func BenchmarkHashURL(b *testing.B) {
	svc := NewURLService(repository.NewURLInMemoryRepository(), "http://localhost:8080", worker.NewDeleteURLWorker(repository.NewURLInMemoryRepository(), 10), audit.NewNoop())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = svc.hashURL("https://example.com/some/really/long/url/path?with=query&and=values")
	}
}

func BenchmarkShortenURL_InMemory(b *testing.B) {
	repo := repository.NewURLInMemoryRepository()
	w := worker.NewDeleteURLWorker(repo, 10)
	svc := NewURLService(repo, "http://localhost:8080", w, audit.NewNoop())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = svc.ShortenURL("https://example.com/some/really/long/url/path?with=query&and=values", "user-1")
	}
}

func BenchmarkBatchShortenURL_InMemory(b *testing.B) {
	repo := repository.NewURLInMemoryRepository()
	w := worker.NewDeleteURLWorker(repo, 10)
	svc := NewURLService(repo, "http://localhost:8080", w, audit.NewNoop())

	items := make([]model.BatchShortenURLRequest, 0, 100)
	for i := 0; i < 100; i++ {
		id := "id"
		items = append(items, model.BatchShortenURLRequest{
			CorrelationID: &id,
			OriginalURL:   "https://example.com/some/really/long/url/path?with=query&and=values",
		})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = svc.BatchShortenURL(items, "user-1")
	}
}
