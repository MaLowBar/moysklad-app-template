package jsonapi

type Project struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func (p Project) GetName() string {
	return p.Name
}

func (p Project) GetID() string {
	return p.Id
}
