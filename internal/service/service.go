package service

import (
	"fmt"
	"io"
	"net/http"
)

type GoogleDocsService struct{}

func NewGoogleDocsService() *GoogleDocsService {
	return &GoogleDocsService{}
}

func (s *GoogleDocsService) DownloadCSV(fileID string) (io.ReadCloser, error) {
	fileURL := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/gviz/tq?tqx=out:csv", fileID)

	resp, err := http.Get(fileURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error downloading CSV file: %w", err)
	}

	return resp.Body, nil
}
