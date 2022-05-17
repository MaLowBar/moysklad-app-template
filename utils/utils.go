package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	HTTPClientTimeout = uint64(60)
	client            = http.Client{Timeout: time.Duration(HTTPClientTimeout) * time.Second}
	Tryes             = uint64(3)
	Timeout           = uint64(3)
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

type MSJsonTime time.Time

func (j *MSJsonTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*j = MSJsonTime(t)
	return nil
}

func (j MSJsonTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}

func Fetch[T any](raw io.ReadCloser) (T, error) {
	var data T
	err := json.NewDecoder(raw).Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func MakeRequest(method, url, accessToken string, body io.ReadCloser) (*http.Response, error) {
	req, err := Request(method, url, accessToken, body)
	if err != nil {
		return nil, err
	}
	var resp *http.Response

	tryCount := Tryes
	for tryCount > 0 {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(Timeout) * time.Second)
		tryCount--
	}
	return resp, err
}
