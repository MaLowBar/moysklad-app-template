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

func (c Cashout) GetName() string {
	return c.Name
}

func (c Cashout) GetID() string {
	return c.Id
}
