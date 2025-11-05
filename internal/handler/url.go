package handler

import (
	"github.com/alikhanturusbekov/go-url-shortener/internal/service"
	"io"
	"log"
	"net/http"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{service: service}
}

func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
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

	url, appError := h.service.ShortenURL(string(body))
	if appError != nil {
		http.Error(w, appError.Message, appError.Code)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(url))
	if err != nil {
		http.Error(w, "failed to write a response", http.StatusBadRequest)
	}
}

func (h *URLHandler) ResolveURL(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	url, appError := h.service.ResolveShortURL(id)
	if appError != nil {
		http.Error(w, appError.Message, appError.Code)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
