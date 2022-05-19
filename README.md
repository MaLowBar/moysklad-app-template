# Микрофреймворк для Go-приложений маркетплейса [МойСклад](https://www.moysklad.ru/ "Ссылка на главную страницу МойСклад")

Данный пакет призван упростить и ускорить написание приложений для маркетплейса системы МойСклад на языке программирования Go. 
На данном этапе реализовано:
* базовое взаимодействие с VendorAPI, т.е. активация, удаление и получение статуса приложения на аккаунте,
* запросы к JSON API 1.2 и несколько конкретных типов сущностей,
* для реализации хранилища данных разработан пакет, который работает с: файлами на диске или СУБД PostgreSQL,
* для работы с шаблонами интерфейса приложения реализован Renderer в ```app.go``` на основе стандартной библиотеки ```html/template```.

Pull requests и Issues приветствуются!

## Описание и пример использования:
Конфигурация приложения создается в виде экземпляра структуры ```moyskladapptemplate.AppConfig```. Далее определяется хранилище данных и его тип. Далее для каждого эндпоинта создается обработчик в виде экземпляра структуры ```moyskladapptemplate.AppHandler``` и все созданные обработчики передаются в метод ```NewApp``` для создания цельного приложения. После этого приложение запускается. Пример кода:
```go
package main

import (
	"fmt"
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
	}

	// Можно использовать БД PostgreSQL
	//myStorage, err := storage.NewPostgreStorage("postgres://msgo:pswd@localhost/msgo_db")
	//if err != nil {
	//	log.Fatal(fmt.Errorf("cannot create storage: %w", err))
	//	return
	//}

	// Инициализируем файловое хранилище
	myStorage, err := storage.NewFileStorage("./")
	if err != nil {
		log.Fatalf("Cannot create app storage: %s", err.Error())
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
			return c.HTML(200, fmt.Sprintf(`<html>
    <head>
    </head>
    <body>
        <center>
            <h1> Hello, %s! </h1>
			<h2> Your id: %s </h2>
			<form action="/go-apps/test-app/get-purchaseorders" method="POST">
			<input type="hidden" name="accountId" value="%s"/>
  			<p><input type="submit" value="Click here"></p>
 			</form> 
        </center>    
    </body>
</html>
`, userContext.FullName, userContext.ID, userContext.AccountID))
		},
	}

	formHandler := moyskladapptemplate.AppHandler{
		Method: "POST",
		Path:   "/go-apps/test-app/get-purchaseorders",
		HandlerFunc: func(c echo.Context) error {
			counterparties, err := jsonapi.GetAllEntities[jsonapi.Counterparty](myStorage, c.FormValue("accountId"), "counterparty")
			if err != nil {
				return &echo.HTTPError{
					Code:    http.StatusInternalServerError,
					Message: err,
				}
			}
			return c.JSON(200, counterparties)
		},
	}
	// Создаем приложение
	app := moyskladapptemplate.NewApp(&info, myStorage, iframeHandler, formHandler)

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
```
