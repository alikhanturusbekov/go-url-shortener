package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HTTPObserver struct {
	client *http.Client
	url    string
}

func NewHTTPObserver(url string) *HTTPObserver {
	return &HTTPObserver{
		client: &http.Client{Timeout: 5 * time.Second},
		url:    url,
	}
}

func (h *HTTPObserver) Send(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	resp, err := h.client.Post(h.url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	return nil
}
