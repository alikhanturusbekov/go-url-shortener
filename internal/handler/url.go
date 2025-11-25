package handler

import (
	"encoding/json"
	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
	"github.com/alikhanturusbekov/go-url-shortener/internal/service"
	"github.com/alikhanturusbekov/go-url-shortener/pkg/logger"
	"go.uber.org/zap"
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

func (h *URLHandler) ShortenURLAsText(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, appError.GetFullMessage(), appError.Code)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(url))
	if err != nil {
		http.Error(w, "failed to write a response", http.StatusBadRequest)
	}
}

func (h *URLHandler) ShortenURLAsJSON(w http.ResponseWriter, r *http.Request) {
	var req model.Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Error("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	url, appError := h.service.ShortenURL(req.URL)
	if appError != nil {
		http.Error(w, appError.GetFullMessage(), appError.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := model.Response{Result: url}

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		http.Error(w, "failed to write a response", http.StatusBadRequest)
		return
	}
}

func (h *URLHandler) ResolveURL(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	url, appError := h.service.ResolveShortURL(id)
	if appError != nil {
		http.Error(w, appError.GetFullMessage(), appError.Code)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
