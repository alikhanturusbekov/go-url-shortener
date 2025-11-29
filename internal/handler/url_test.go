package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alikhanturusbekov/go-url-shortener/internal/config"
	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
	"github.com/alikhanturusbekov/go-url-shortener/internal/repository"
	"github.com/alikhanturusbekov/go-url-shortener/internal/service"
)

var testConfig *config.Config

var database *sql.DB

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "url_test")
	if err != nil {
		fmt.Println("Failed to create temp directory:", err)
		os.Exit(1)
	}

	testConfig = &config.Config{
		Address:         "localhost:9999",
		BaseURL:         "http://localhost:9999",
		FileStoragePath: filepath.Join(tmpDir, "url_pairs_test.json"),
		DatabaseDSN:     "postgres://username:password@localhost:5432/shortened_urls",
	}

	database, err := sql.Open("pgx", testConfig.DatabaseDSN)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	code := m.Run()

	err = os.RemoveAll(tmpDir)
	if err != nil {
		fmt.Println("Failed to delete temp directory:", err)
		os.Exit(1)
	}

	os.Exit(code)
}

func TestShortenURLAsText(t *testing.T) {
	type requestData struct {
		headers map[string]string
		method  string
		body    io.Reader
	}
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name        string
		requestData requestData
		want        want
	}{
		{
			name: "Positive case: URL with scheme",
			requestData: requestData{
				headers: map[string]string{"Content-Type": "text/plain"},
				method:  http.MethodPost,
				body:    strings.NewReader("https://yandex.ru"),
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
			},
		},
		{
			name: "Negative case: URL without scheme",
			requestData: requestData{
				headers: map[string]string{"Content-Type": "text/plain"},
				method:  http.MethodPost,
				body:    strings.NewReader("yandex.ru"),
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name: "Negative Case: Empty URL",
			requestData: requestData{
				headers: map[string]string{"Content-Type": "text/plain"},
				method:  http.MethodPost,
				body:    strings.NewReader(""),
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.requestData.method, "/", tt.requestData.body)
			for name, value := range tt.requestData.headers {
				request.Header.Set(name, value)
			}

			urlRepo, err := setupURLFileRepository(testConfig.FileStoragePath)
			require.NoError(t, err)
			urlService := service.NewURLService(urlRepo, testConfig.BaseURL)
			h := NewURLHandler(urlService, database).ShortenURLAsText
			w := httptest.NewRecorder()
			h(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if tt.want.statusCode == http.StatusCreated {
				resultBody, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)
				require.NotEmpty(t, resultBody)

				resultURL, err := url.Parse(string(resultBody))
				require.NoError(t, err)
				assert.Equal(t, "http", resultURL.Scheme)

				shortenedPath := strings.TrimPrefix(resultURL.Path, "/")
				assert.Equal(t, 7, len(shortenedPath))
			}
		})
	}
}

func TestShortenURLAsJSON(t *testing.T) {
	type requestData struct {
		headers map[string]string
		method  string
		body    []byte
	}
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name        string
		requestData requestData
		want        want
	}{
		{
			name: "Positive case: URL with scheme",
			requestData: requestData{
				headers: map[string]string{"Content-Type": "application/json"},
				method:  http.MethodPost,
				body:    []byte(`{"url": "https://yandex.ru"}`),
			},
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
			},
		},
		{
			name: "Negative case: incorrect URL field",
			requestData: requestData{
				headers: map[string]string{"Content-Type": "application/json"},
				method:  http.MethodPost,
				body:    []byte(`{"incorrectURL": "https://yandex.ru"}`),
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.requestData.method, "/", bytes.NewReader(tt.requestData.body))
			for name, value := range tt.requestData.headers {
				request.Header.Set(name, value)
			}

			urlRepo, err := setupURLFileRepository(testConfig.FileStoragePath)
			require.NoError(t, err)
			urlService := service.NewURLService(urlRepo, testConfig.BaseURL)
			h := NewURLHandler(urlService, database).ShortenURLAsJSON
			w := httptest.NewRecorder()
			h(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			fmt.Println(result.Body)

			if tt.want.statusCode == http.StatusCreated {
				var response model.Response
				err := json.NewDecoder(result.Body).Decode(&response)
				require.NoError(t, err)

				resultURL, err := url.Parse(response.Result)
				require.NoError(t, err)
				assert.Equal(t, "http", resultURL.Scheme)

				shortenedPath := strings.TrimPrefix(resultURL.Path, "/")
				assert.Equal(t, 7, len(shortenedPath))
			}
		})
	}
}

func TestResolveURL(t *testing.T) {
	type want struct {
		statusCode int
		headers    map[string]string
	}
	tests := []struct {
		name            string
		targetURL       string
		mockURLDatabase []model.URLPair
		want            want
	}{
		{
			name:      "Positive case: URL exists in Database",
			targetURL: "abcdefg",
			mockURLDatabase: []model.URLPair{
				{
					Short: "abcdefg",
					Long:  "https://yandex.ru",
				},
			},
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				headers: map[string]string{
					"Content-Type": "text/html; charset=utf-8",
					"Location":     "https://yandex.ru",
				},
			},
		},
		{
			name:      "Negative case: URL does not exist in Database",
			targetURL: "abcdefg",
			mockURLDatabase: []model.URLPair{
				{
					Short: "gfedvba",
					Long:  "https://yandex.ru",
				},
			},
			want: want{
				statusCode: http.StatusNotFound,
				headers:    map[string]string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlRepo, err := setupURLFileRepository(testConfig.FileStoragePath)
			require.NoError(t, err)

			for _, urlPair := range tt.mockURLDatabase {
				err := urlRepo.Save(urlPair)
				require.NoError(t, err)
			}

			mux := http.NewServeMux()
			urlService := service.NewURLService(urlRepo, testConfig.BaseURL)
			mux.HandleFunc("/{id}", NewURLHandler(urlService, database).ResolveURL)

			request := httptest.NewRequest(http.MethodGet, "/"+tt.targetURL, nil)

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			for header, value := range tt.want.headers {
				assert.Equal(t, value, result.Header.Get(header))
			}
		})
	}
}

func setupURLFileRepository(filePath string) (*repository.URLFileRepository, error) {
	err := os.WriteFile(filePath, []byte(""), 0644)
	if err != nil {
		return nil, err
	}

	urlRepo, err := repository.NewURLFileRepository(filePath)
	if err != nil {
		return nil, err
	}

	return urlRepo, nil
}
