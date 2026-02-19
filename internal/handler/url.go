package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
	"github.com/alikhanturusbekov/go-url-shortener/internal/service"
	"github.com/alikhanturusbekov/go-url-shortener/pkg/authorization"
	"github.com/alikhanturusbekov/go-url-shortener/pkg/logger"
)

type URLHandler struct {
	service  *service.URLService
	database *sql.DB
}

func NewURLHandler(service *service.URLService, database *sql.DB) *URLHandler {
	return &URLHandler{
		service:  service,
		database: database,
	}
}

func (h *URLHandler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.database.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	url, appError := h.service.ShortenURL(string(body), h.getUserID(r))
	if appError != nil && url == "" {
		http.Error(w, appError.GetFullMessage(), appError.Code)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if appError != nil {
		w.WriteHeader(appError.Code)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

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

	url, appError := h.service.ShortenURL(req.URL, h.getUserID(r))
	if appError != nil && url == "" {
		http.Error(w, appError.GetFullMessage(), appError.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if appError != nil {
		w.WriteHeader(appError.Code)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

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

func (h *URLHandler) BatchShortenURL(w http.ResponseWriter, r *http.Request) {
	var req []model.BatchShortenURLRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Error("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	results, err := h.service.BatchShortenURL(req, h.getUserID(r))
	if err != nil {
		http.Error(w, "failed to batch shorten", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *URLHandler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := authorization.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "need to authorize to access this method", http.StatusUnauthorized)
	}

	userURLs, err := h.service.GetUserURLs(userID)
	if err != nil {
		http.Error(w, "failed to fetch user URLs", http.StatusInternalServerError)
	}

	if len(userURLs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(userURLs); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *URLHandler) DeleteUserURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := authorization.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "need to authorize to access this method", http.StatusUnauthorized)
	}

	var shorts []string
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&shorts); err != nil {
		logger.Log.Error("cannot decode request JSON body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteUserURLs(userID, shorts)
	if err != nil {
		http.Error(w, "failed to delete URL pairs", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *URLHandler) getUserID(r *http.Request) string {
	userID, ok := authorization.UserIDFromContext(r.Context())
	if !ok {
		return ""
	}

	return userID
}
