package storage

import (
	"encoding/json"
	"fmt"
	moyskladapptemplate "github.com/MaLowBar/moysklad-app-template"
	"os"
	"strings"
)

type FileStorage struct {
	path string
	apps map[string]appInfo
}

// NewFileStorage returns new FileStorage with configured path. Path must have "/" postfix.
func NewFileStorage(path string) (*FileStorage, error) {
	apps := make(map[string]appInfo)
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".app") {
			var data []byte
			data, err = os.ReadFile(file.Name())
			if err != nil {
				return nil, err
			}
			var app appInfo
			if err = json.Unmarshal(data, &app); err != nil {
				return nil, err
			}
			apps[app.AccountId] = app
		}
	}
	return &FileStorage{path: path, apps: apps}, nil
}

type appInfo struct {
	AccountId   string                        `json:"account_id"`
	Status      moyskladapptemplate.AppStatus `json:"status"`
	AccessToken string                        `json:"access_token"`
}

func (fs *FileStorage) Activate(accountId, accessToken string) (moyskladapptemplate.AppStatus, error) {
	app := appInfo{AccountId: accountId, Status: moyskladapptemplate.StatusActivated, AccessToken: accessToken}
	data, err := json.Marshal(app)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(fmt.Sprintf("%s%s.app", fs.path, accountId), data, 0644)
	if err != nil {
		return "", err
	}

	fs.apps[accountId] = app

	return app.Status, nil
}

func (fs *FileStorage) Delete(accountId string) error {
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

	delete(fs.apps, accountId)

	return nil
}

func (fs *FileStorage) GetStatus(accountId string) (moyskladapptemplate.AppStatus, error) {
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

func (fs *FileStorage) AccessTokenByAccountId(accountId string) (string, error) {
	if a, ok := fs.apps[accountId]; ok {
		return a.AccessToken, nil
	}
	return "", fmt.Errorf("no app associated with this account id: %s", accountId)
}
