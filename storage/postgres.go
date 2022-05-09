package storage

import (
	"database/sql"
	"fmt"
	templ "github.com/MaLowBar/moysklad-app-template"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgreStorage struct {
	db   *sql.DB
	apps map[string]appInfo
}

func NewPostgreStorage(connect string) (*PostgreStorage, error) {
	db, err := sql.Open("pgx", connect)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS apps (accountId varchar, status varchar, accessToken varchar)`)
	if err != nil {
		return nil, err
	}

	apps := make(map[string]appInfo)
	rows, err := db.Query(`SELECT accountId, status, accessToken FROM apps`)
	defer rows.Close()

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var app appInfo
		err = rows.Scan(&app.AccountId, &app.Status, &app.AccessToken)
		if err != nil {
			return nil, err
		}
		apps[app.AccountId] = app
	}
	return &PostgreStorage{db: db, apps: apps}, nil
}

func (s *PostgreStorage) Activate(accountId, accessToken string) (templ.AppStatus, error) {
	_, err := s.db.Exec(`INSERT INTO apps VALUES ($1, $2, $3)`, accountId, templ.StatusActivated, accessToken)
	if err != nil {
		return "", err
	}

	app := appInfo{AccountId: accountId, Status: templ.StatusActivated, AccessToken: accessToken}
	s.apps[accountId] = app

	return templ.StatusActivated, nil
}

func (s *PostgreStorage) Delete(accountId string) error {
	_, err := s.db.Exec(`DELETE FROM apps WHERE accountId = $1`, accountId)
	if err != nil {
		return err
	}

	delete(s.apps, accountId)

	return nil
}

func (s *PostgreStorage) GetStatus(accountId string) (templ.AppStatus, error) {
	row := s.db.QueryRow(`SELECT status FROM apps WHERE accountId = $1`, accountId)
	if err := row.Err(); err != nil {
		return "", err
	}
	var status templ.AppStatus
	if err := row.Scan(&status); err != nil {
		return "", err
	}
	return status, nil
}

func (s *PostgreStorage) AccessTokenByAccountId(accountId string) (string, error) {
	if a, ok := s.apps[accountId]; ok {
		return a.AccessToken, nil
	}
	return "", fmt.Errorf("no app associated with this account id: %s", accountId)
}
