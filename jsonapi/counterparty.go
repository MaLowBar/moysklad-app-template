package jsonapi

type Counterparty struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func (c Counterparty) GetName() string {
	return c.Name
}

func (c Counterparty) GetID() string {
	return c.Id
}
