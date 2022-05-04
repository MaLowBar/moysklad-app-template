package moyskladapptemplate

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

type AppConfig struct {
	ID           string
	UID          string
	SecretKey    string
	AccessToken  string
	VendorAPIURL string
}

type AppStatus string

const (
	StatusActivated        AppStatus = "Activated"
	StatusSettingsRequired AppStatus = "Settings required"
	StatusActivating       AppStatus = "Activating"
	StatusSuspended        AppStatus = "Suspended"
)

type AppStorage interface {
	Activate(accountId, accessToken string) (AppStatus, error)
	Delete(accountId string) error
	GetStatus(accountId string) (AppStatus, error)
	AccessTokenByAccountId(accountId string) (string, error)
}

type AppHandler struct {
	Method string
	Path   string
	echo.HandlerFunc
}

type App struct {
	info    *AppConfig
	storage AppStorage
	srv     *echo.Echo
}

func NewApp(appConfig *AppConfig, storage AppStorage, handlers ...AppHandler) *App {
	app := &App{
		info:    appConfig,
		storage: storage,
	}

	srv := echo.New()

	srv.Use(middleware.Logger(), middleware.Recover())

	vendorAPIURL := appConfig.VendorAPIURL
	srv.Add("PUT", vendorAPIURL, app.activateHandler)
	srv.Add("DELETE", vendorAPIURL, app.deleteHandler)
	srv.Add("GET", vendorAPIURL, app.getStatusHandler)

	for _, handler := range handlers {
		srv.Add(handler.Method, handler.Path, handler.HandlerFunc)
	}

	app.srv = srv

	return app
}

func (a *App) Run(addr string) error {
	return a.srv.Start(addr)
}

func (a *App) Stop(timeout int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	return a.srv.Shutdown(ctx)
}

type activateReq struct {
	Access []struct {
		AccessToken string `json:"access_token"`
	} `json:"access,omitempty"`
}

func (a *App) activateHandler(c echo.Context) error {
	if a.info.ID != c.Param("appId") {
		return c.NoContent(http.StatusNotFound)
	}

	var req activateReq
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	if status, err := a.storage.Activate(c.Param("accountId"), req.Access[0].AccessToken); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	} else {
		a.info.AccessToken = req.Access[0].AccessToken
		return c.JSON(http.StatusOK, map[string]string{"status": string(status)})
	}
}

func (a *App) deleteHandler(c echo.Context) error {
	if a.info.ID != c.Param("appId") {
		return c.NoContent(http.StatusNotFound)
	}

	if err := a.storage.Delete(c.Param("accountId")); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.NoContent(http.StatusOK)
}

func (a *App) getStatusHandler(c echo.Context) error {
	if a.info.ID != c.Param("appId") {
		return c.NoContent(http.StatusNotFound)
	}

	status, err := a.storage.GetStatus(c.Param("accountId"))
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	if status == StatusSuspended {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, map[string]string{"status": string(status)})
}
