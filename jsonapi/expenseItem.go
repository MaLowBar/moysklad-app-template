package jsonapi

type ExpenseItem struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type ExpenseItemForReq struct {
	ExpenseItem struct {
		Meta struct {
			Href         string `json:"href"`
			MetaDataHref string `json:"metadataHref"`
			Type         string `json:"type"`
			MediaType    string `json:"mediaType"`
		} `json:"meta"`
	} `json:"expenseItem"`
}

func (e ExpenseItem) GetName() string {
	return e.Name
}

func (e ExpenseItem) GetID() string {
	return e.Id
}
