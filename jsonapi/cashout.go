package jsonapi

type Cashout struct {
	Name  string `json:"name"`
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
	Description string `json:"description"`
}
