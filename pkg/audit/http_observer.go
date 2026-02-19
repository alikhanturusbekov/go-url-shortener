package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HTTPObserver structure to observe audit events and sends them to audit service
type HTTPObserver struct {
	client *http.Client
	url    string
}

// NewHTTPObserver creates a new HTTPObserver
func NewHTTPObserver(url string) *HTTPObserver {
	return &HTTPObserver{
		client: &http.Client{Timeout: 5 * time.Second},
		url:    url,
	}
}

// Send sends event logs to audit service
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
