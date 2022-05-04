package storage

import (
	"encoding/json"
	"fmt"
	moyskladapptemplate "github.com/MaLowBar/moysklad-app-template"
	"os"
)

type FileStorage struct {
	*moyskladapptemplate.AppConfig
	path string
}

// NewFileStorage returns new FileStorage with configured path. Path must have "/" postfix.
func NewFileStorage(info *moyskladapptemplate.AppConfig, path string) FileStorage {
	return FileStorage{AppConfig: info, path: path}
}

type appInfo struct {
	AccountId   string                        `json:"account_id"`
	Status      moyskladapptemplate.AppStatus `json:"status"`
	AccessToken string                        `json:"access_token"`
}

func (fs FileStorage) Activate(accountId, accessToken string) (moyskladapptemplate.AppStatus, error) {
	app := appInfo{AccountId: accountId, Status: moyskladapptemplate.StatusActivated, AccessToken: accessToken}
	data, err := json.Marshal(app)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(fmt.Sprintf("%s%s.app", fs.path, accountId), data, 0644)
	if err != nil {
		return "", err
	}
	fs.AppConfig.AccessToken = accessToken
	return app.Status, nil
}

func (fs FileStorage) Delete(accountId string) error {
	data, err := os.ReadFile(fmt.Sprintf("%s%s.app", fs.path, accountId))
	if err != nil {
		return err
	}
	var app appInfo
	if err = json.Unmarshal(data, &app); err != nil {
		return err
	}
	app.Status = moyskladapptemplate.StatusSuspended
	data, err = json.Marshal(app)
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s%s.app", fs.path, accountId), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (fs FileStorage) GetStatus(accountId string) (moyskladapptemplate.AppStatus, error) {
	data, err := os.ReadFile(fmt.Sprintf("%s%s.app", fs.path, accountId))
	if err != nil {
		return "", err
	}
	var app appInfo
	if err = json.Unmarshal(data, &app); err != nil {
		return "", err
	}
	return app.Status, nil
}
