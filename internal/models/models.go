package models

import "fmt"

type BatchInput struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchOutput struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type BatchForPost struct {
	CorrelationID string
	OriginalURL   string
	ShortURL      string
}

func (b *BatchInput) Validate() bool {
	if b.OriginalURL != "" {
		return true
	}
	return false
}

type ArrBatchInput []BatchInput

func (ab ArrBatchInput) Validate() bool {
	result := false
	if len(ab) == 0 {
		return false
	}
	for _, val := range ab {
		result = result || val.Validate()
	}
	return result
}

type DBError struct {
	LongURL string
	Err     error
}

func (dbe *DBError) Error() string {
	return fmt.Sprintf("[%s] %v", dbe.LongURL, dbe.Err)
}

func NewDBError(longUrl string, err error) error {
	return &DBError{LongURL: longUrl, Err: err}
}
