package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/alikhanturusbekov/go-url-shortener/internal/repository"
	"github.com/alikhanturusbekov/go-url-shortener/internal/service"
	"github.com/alikhanturusbekov/go-url-shortener/internal/worker"
	"github.com/alikhanturusbekov/go-url-shortener/pkg/audit"
)

// BatchRequest represents batch shorten request
type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchResponse represents batch shorten response
type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// Example demonstrates HTTP usage of URLHandler endpoints.
func Example() {
	r := chi.NewRouter()

	repo := repository.NewURLInMemoryRepository()
	deleteWorker := worker.NewDeleteURLWorker(repo, 10)
	urlService := service.NewURLService(repo, "http://localhost:8080", deleteWorker, audit.NewNoop())
	handler := NewURLHandler(urlService, nil)

	r.Post("/api/shorten", handler.ShortenURLAsJSON)
	r.Get("/{id}", handler.ResolveURL)
	r.Get("/api/user/urls", handler.GetUserURLs)
	r.Post("/api/shorten/batch", handler.BatchShortenURL)

	ts := httptest.NewServer(r)
	defer ts.Close()

	// 1. Shorten URL (JSON)

	jsonBody := `{"url":"https://google.com"}`
	resp, _ := http.Post(ts.URL+"/api/shorten", "application/json", strings.NewReader(jsonBody))
	fmt.Printf("POST /api/shorten: %d\n", resp.StatusCode)
	resp.Body.Close()

	// 2. Resolve short URL

	short, _ := urlService.ShortenURL("https://google.com", "")
	shortID := short[len("http://localhost:8080/"):]

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, _ = client.Get(ts.URL + "/" + shortID)
	fmt.Printf("GET /{id}: %d\n", resp.StatusCode)
	resp.Body.Close()

	// 3. User URLs

	resp, _ = client.Get(ts.URL + "/api/user/urls")
	fmt.Printf("GET /api/user/urls: %d\n", resp.StatusCode)
	resp.Body.Close()

	// 4. Batch shorten

	input := []BatchRequest{
		{CorrelationID: "1", OriginalURL: "https://google.com"},
		{CorrelationID: "2", OriginalURL: "https://yandex.ru"},
	}

	body, _ := json.Marshal(input)

	resp, _ = client.Post(ts.URL+"/api/shorten/batch", "application/json", bytes.NewBuffer(body))
	fmt.Printf("POST /api/shorten/batch: %d\n", resp.StatusCode)
	defer resp.Body.Close()

	var result []BatchResponse
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Printf("Batch items: %d\n", len(result))

	// Output:
	// POST /api/shorten: 201
	// GET /{id}: 307
	// GET /api/user/urls: 401
	// POST /api/shorten/batch: 201
	// Batch items: 2
}
