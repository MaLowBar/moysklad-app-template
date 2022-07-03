package storage

import moyskladapptemplate "github.com/MaLowBar/moysklad-app-template"

type AppInfo struct {
	AccountId   string                        `json:"account_id"`
	Status      moyskladapptemplate.AppStatus `json:"status"`
	AccessToken string                        `json:"access_token"`
}
