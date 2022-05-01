package moyskladapptemplate

import "fmt"

type JSONAPIError struct {
	Header       string `json:"error"`
	Parameter    string `json:"parameter,omitempty"`
	Code         int    `json:"code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

func (e JSONAPIError) Error() string {
	return fmt.Sprintf("Error: %s, Parameter: %s, Code: %d, ErrorMessage: %s",
		e.Header, e.Parameter, e.Code, e.ErrorMessage)
}
