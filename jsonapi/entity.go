package jsonapi

type Entity struct {
	Id    string `json:"id"`
	Agent struct {
		Meta struct {
			Href string `json:"href"`
		} `json:"meta"`
	} `json:"agent"`
	Project struct {
		Meta struct {
			Href string `json:"href"`
		} `json:"meta"`
	} `json:"project"`
	Description    string `json:"description"`
	PaymentPurpose string `json:"paymentPurpose"`
}
