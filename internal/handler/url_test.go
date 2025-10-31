package handler

import (
	"github.com/alikhanturusbekov/go-url-shortener/internal/repository"
	"github.com/alikhanturusbekov/go-url-shortener/internal/service"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenURL(t *testing.T) {
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
			name: "Positive case: URL without scheme",
			requestData: requestData{
				headers: map[string]string{"Content-Type": "text/plain"},
				method:  http.MethodPost,
				body:    strings.NewReader("yandex.ru"),
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
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
		{
			name: "Negative Case: Sent Content-Type application/json",
			requestData: requestData{
				headers: map[string]string{"Content-Type": "application/json"},
				method:  http.MethodPost,
				body:    strings.NewReader("https://yandex.ru"),
			},
			want: want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name: "Negative Case: GET request",
			requestData: requestData{
				headers: map[string]string{"Content-Type": "text/plain"},
				method:  http.MethodGet,
				body:    strings.NewReader("https://yandex.ru"),
			},
			want: want{
				contentType: "",
				statusCode:  http.StatusMethodNotAllowed,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.requestData.method, "/", tt.requestData.body)
			for name, value := range tt.requestData.headers {
				request.Header.Set(name, value)
			}

			urlRepo := repository.NewURLMapRepository()
			urlService := service.NewURLService(urlRepo)
			h := NewURLHandler(urlService).ShortenURL
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

func TestResolveURL(t *testing.T) {
	type want struct {
		statusCode int
		headers    map[string]string
	}
	tests := []struct {
		name            string
		targetURL       string
		mockURLDatabase map[string]string
		want            want
	}{
		{
			name:      "Positive case: URL exists in Database",
			targetURL: "abcdefg",
			mockURLDatabase: map[string]string{
				"abcdefg": "https://yandex.ru",
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
			mockURLDatabase: map[string]string{
				"gfedvba": "https://yandex.ru",
			},
			want: want{
				statusCode: http.StatusNotFound,
				headers:    map[string]string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlRepo := repository.NewURLMapRepository()

			for short, long := range tt.mockURLDatabase {
				err := urlRepo.Save(short, long)
				require.NoError(t, err)
			}

			mux := http.NewServeMux()
			urlService := service.NewURLService(urlRepo)
			mux.HandleFunc("/{id}", NewURLHandler(urlService).ResolveURL)

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
