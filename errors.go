package moyskladapptemplate

import "fmt"

type JSONAPIError struct {
	Errors []struct {
		Header       string `json:"error"`
		Parameter    string `json:"parameter,omitempty"`
		Code         int    `json:"code,omitempty"`
		ErrorMessage string `json:"error_message,omitempty"`
	} `json:"errors"`
}

func (e JSONAPIError) Error() (res string) {
	for _, er := range e.Errors {
		res += fmt.Sprintf("Error: %s, Parameter: %s, Code: %d, ErrorMessage: %s",
			er.Header, er.Parameter, er.Code, er.ErrorMessage)
	}
	return
}
