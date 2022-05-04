package main

import (
	"fmt"
	"github.com/MaLowBar/moysklad-app-template"
	"github.com/MaLowBar/moysklad-app-template/jsonapi"
	"github.com/MaLowBar/moysklad-app-template/storage"
	"github.com/MaLowBar/moysklad-app-template/vendorapi"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	// Задаём информацию необходимую для работы приложения
	var info = moyskladapptemplate.AppConfig{
		ID:           "fb08e3f3-8f1a-488e-a609-1baa389cc546",
		UID:          "purchases-report-go.sorochinsky",
		SecretKey:    "8iv6RbvFlQsiDMqQz4ECczLjiwEZRfBkVKa2cMBmsHnzIg2ELuqdbQNXvloY65nQD1crmxdbCVXbx1CvnjY1Th9sUebNXOYnULPtZ40N2ujjv7EzbE6F5SEM9xucnEAL",
		VendorAPIURL: "/echo/api/moysklad/vendor/1.0/apps/:appId/:accountId",
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
		Path:   "/echo/iframe/purchases-report-go.sorochinsky",
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
			<form action="/echo/test-get-purchaseorders" method="POST">
  			<p><input type="submit" value="Click here"></p>
 			</form> 
        </center>    
    </body>
</html>
`, userContext.FullName, userContext.ID))
		},
	}

	formHandler := moyskladapptemplate.AppHandler{
		Method: "POST",
		Path:   "/echo/test-get-purchaseorders",
		HandlerFunc: func(c echo.Context) error {
			orders, err := jsonapi.GetPurchaseOrders(info)
			if err != nil {
				return &echo.HTTPError{
					Code:    http.StatusInternalServerError,
					Message: err,
				}
			}
			return c.JSON(200, orders)
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
	case err := <-e:
		log.Printf("Server returned error: %s", err)
	case <-c:
		app.Stop(5)
		log.Println("Stop signal received")
	}
}
