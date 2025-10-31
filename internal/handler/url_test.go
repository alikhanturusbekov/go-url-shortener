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
			for k, v := range tt.requestData.headers {
				request.Header.Set(k, v)
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
