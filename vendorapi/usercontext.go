package vendorapi

// UserContext представляет собой данные о сотруднике, от лица которого происходит запрос. Подробнее смотри https://dev.moysklad.ru/doc/api/remap/1.2/#mojsklad-json-api-obschie-swedeniq-kontext-zaprosa-sotrudnika
type UserContext struct {
	ID          string `json:"id"`
	UID         string `json:"uid,omitempty"`
	AccountID   string `json:"accountId"`
	FullName    string `json:"fullName,omitempty"`
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Position    string `json:"position,omitempty"`
	Permissions []struct {
		Admin struct {
			View string `json:"view"`
		} `json:"admin,omitempty"`
	} `json:"permissions,omitempty"`
}
