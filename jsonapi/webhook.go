package jsonapi

type WebhookResponse struct {
	Events []struct {
		Meta struct {
			EssenceType string `json:"type"`
			Href        string `json:"href"`
		} `json:"meta"`
		Action    string `json:"action"`
		AccountId string `json:"accountId"`
	} `json:"events"`
}
