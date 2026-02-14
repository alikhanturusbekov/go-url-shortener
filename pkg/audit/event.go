package audit

type Event struct {
	TS     int64  `json:"ts"`
	Action string `json:"action"` // shorten | follow
	UserID string `json:"user_id,omitempty"`
	URL    string `json:"url"`
}
