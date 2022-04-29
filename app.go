package moyskladapptemplate

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type AppInfo struct {
	ID        string
	SecretKey string
}

type AppStatus string

const (
	StatusActivated        AppStatus = "Activated"
	StatusSettingsRequired AppStatus = "Settings required"
	StatusActivating       AppStatus = "Activating"
	StatusInactive         AppStatus = "Inactive"
)

type AppStorage interface {
	Activate(accountId, accessToken string) (AppStatus, error)
	Delete(accountId string) error
	GetStatus(accountId string) (AppStatus, error)
}

type AppHandler struct {
	Method string
	Path   string
	echo.HandlerFunc
}

type App struct {
	info    AppInfo
	storage AppStorage
	srv     *echo.Echo
}

const vendorAPIURL = "/api/moysklad/vendor/1.0/apps/:appId/:accountId"

func NewApp(appInfo AppInfo, storage AppStorage, handlers ...AppHandler) *App {
	app := &App{
		info:    appInfo,
		storage: storage,
	}

	srv := echo.New()

	srv.Use(middleware.Logger(), middleware.Recover(), middleware.JWT(appInfo.SecretKey))

	srv.Add("PUT", vendorAPIURL, app.activateHandler)
	srv.Add("DELETE", vendorAPIURL, app.deleteHandler)
	srv.Add("GET", vendorAPIURL, app.getStatusHandler)

	for _, handler := range handlers {
		srv.Add(handler.Method, handler.Path, handler.HandlerFunc)
	}

	app.srv = srv

	return app
}

func (a App) Run(addr string) error {
	return a.srv.Start(addr)
}

func (a App) Stop() error {
	return a.srv.Shutdown(context.Background())
}

type activateReq struct {
	Access []struct {
		AccessToken string `json:"access_token"`
	} `json:"access,omitempty"`
}

func (a App) activateHandler(c echo.Context) error {
	if a.info.ID != c.Param("appId") {
		return c.NoContent(http.StatusNotFound)
	}

	var req activateReq
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if status, err := a.storage.Activate(c.Param("accountId"), req.Access[0].AccessToken); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	} else {
		return c.JSON(http.StatusOK, map[string]string{"status": string(status)})
	}
}

func (a App) deleteHandler(c echo.Context) error {
	if a.info.ID != c.Param("appId") {
		return c.NoContent(http.StatusNotFound)
	}

	if err := a.storage.Delete(c.Param("accountId")); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (a App) getStatusHandler(c echo.Context) error {
	if a.info.ID != c.Param("appId") {
		return c.NoContent(http.StatusNotFound)
	}

	if status, err := a.storage.GetStatus(c.Param("accountId")); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	} else {
		return c.JSON(http.StatusOK, map[string]string{"status": string(status)})
	}
}
