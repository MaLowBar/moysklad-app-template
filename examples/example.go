package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	moyskladapptemplate "github.com/MaLowBar/moysklad-app-template"
	"github.com/MaLowBar/moysklad-app-template/jsonapi"
	"github.com/MaLowBar/moysklad-app-template/storage"
	"github.com/MaLowBar/moysklad-app-template/vendorapi"
	"github.com/labstack/echo/v4"
)

func main() {
	// Задаём информацию необходимую для работы приложения
	var info = moyskladapptemplate.AppConfig{
		ID:           "fb08e3f3-8f1a-488e-a609-1baa389cc546",
		UID:          "test-app.sorochinsky",
		SecretKey:    "8iv6RbvFlQsiDMqQz4ECczLjiwEZRfBkVKa2cMBmsHnzIg2ELuqdbQNXvloY65nQD1crmxdbCVXbx1CvnjY1Th9sUebNXOYnULPtZ40N2ujjv7EzbE6F5SEM9xucnEAL",
		VendorAPIURL: "/go-apps/test-app/api/moysklad/vendor/1.0/apps/:appId/:accountId",
		AppURL:       "https://dev1.the-progress-machine.ru/go-apps/dev1", // URL приложения
	}
	//Если есть необходимость в веб-хуках, создаем мапу вида ["entityType"][]string{"ACTION"} и передаем в info.WebHooksMap
	// whMap := make(map[string][]string)
	// whMap["cashout"] = []string{"CREATE", "UPDATE"}
	// info.WebHooksMap = whMap

	// Можно использовать БД PostgreSQL
	//myStorage, err := storage.NewPostgreStorage("postgres://msgo:pswd@localhost/msgo_db")
	//if err != nil {
	//	log.Panicf(fmt.Errorf("cannot create storage: %w", err))
	//	return
	//}

	// Инициализируем файловое хранилище
	myStorage, err := storage.NewFileStorage("./")
	if err != nil {
		log.Panicf("Cannot create app storage: %s", err.Error())
	}

	// Определяем простейший обработчик для HTML-документа
	var iframeHandler = moyskladapptemplate.AppHandler{
		Method: "GET",
		Path:   "/go-apps/test-app/iframe",
		HandlerFunc: func(c echo.Context) error {
			userContext, err := vendorapi.GetUserContext(c.QueryParam("contextKey"), info)
			if err != nil {
				return &echo.HTTPError{
					Code:    http.StatusInternalServerError,
					Message: err,
				}
			}
			return c.Render(http.StatusOK, "iframe", map[string]interface{}{
				"fullName":  userContext.FullName,
				"accountId": userContext.AccountID,
			})
		},
	}

	formHandler := moyskladapptemplate.AppHandler{
		Method: "POST",
		Path:   "/go-apps/test-app/get-counterparties",
		HandlerFunc: func(c echo.Context) error {
			counterparties, err := jsonapi.GetAllEntities[jsonapi.Counterparty](myStorage, c.FormValue("accountId"), "counterparty", "")
			if err != nil {
				return &echo.HTTPError{
					Code:    http.StatusInternalServerError,
					Message: err,
				}
			}
			return c.Render(http.StatusOK, "iframe", map[string]interface{}{
				"successMessage": "Some message!",
				"list":           counterparties,
			})
		},
	}

	//Формируем слайс с именами шаблонов. Например: []string{"header.html", "footer.html"}
	templateNames := []string{}

	// Создаем приложение
	app := moyskladapptemplate.NewApp(&info, myStorage, templateNames, iframeHandler, formHandler)

	e := make(chan error)
	go func() {
		e <- app.Run("0.0.0.0:8002") // Запускаем
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	select {
	case err = <-e:
		log.Printf("Server returned error: %s", err)
	case <-c:
		app.Stop(5)
		log.Println("Stop signal received")
	}
}
