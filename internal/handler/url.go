package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/alikhanturusbekov/go-url-shortener/internal/service"
)

type UrlHandler struct {
	service *service.UrlService
}

func NewUrlHandler(service *service.UrlService) *UrlHandler {
	return &UrlHandler{service: service}
}

func (h *UrlHandler) ShortenUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close request body read: %s", err)
		}
	}(r.Body)

	original := strings.TrimSpace(string(body))
	if original == "" {
		http.Error(w, "empty URL", http.StatusBadRequest)
		return
	}

	url, _ := h.service.ShortenUrl(string(body))
	fullURL := fmt.Sprintf("%s://%s/%s", "http", r.Host, url)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(fullURL))
	if err != nil {
		http.Error(w, "failed to write a response", http.StatusBadRequest)
	}
}

func (h *UrlHandler) ResolveUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	id := r.PathValue("id")
	url, _ := h.service.ResolveShortUrl(id)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
