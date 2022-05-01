package utils

import (
	"io"
	"net/http"
)

func Request(method, url, bearerToken string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = map[string][]string{
		"Authorization": {"Bearer " + bearerToken},
		"Content-Type":  {"application/json"},
	}
	return req, nil
}
