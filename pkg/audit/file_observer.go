package audit

import (
	"encoding/json"
	"os"
)

// FileObserver structure to observe audit events and writes them to the file
type FileObserver struct {
	file *os.File
}

// NewFileObserver creates new file observer
func NewFileObserver(path string) (*FileObserver, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &FileObserver{file: f}, nil
}

// Send writes the event log to the file
func (f *FileObserver) Send(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = f.file.Write(append(data, '\n'))
	return err
}
