package model

import (
	"github.com/google/uuid"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type URLPair struct {
	ID    string `json:"id"`
	Short string `json:"short"`
	Long  string `json:"long"`
}

func NewURLPair(short, long string, id *string) *URLPair {
	urlPair := &URLPair{
		ID:    uuid.NewString(),
		Short: short,
		Long:  long,
	}

	if id != nil {
		urlPair.ID = *id
	}

	return urlPair
}
