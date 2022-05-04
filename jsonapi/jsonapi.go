package jsonapi

import (
	"encoding/json"
	moyskladapptemplate "github.com/MaLowBar/moysklad-app-template"
	"github.com/MaLowBar/moysklad-app-template/utils"
	"net/http"
	"time"
)

// TODO: сделать получение информации о заказах поставщикам, счетах поставщиков, приемках,
//  заказах покупателей, счетах покупателям, отгрузках,

const jsonEndpoint = "https://online.moysklad.ru/api/remap/1.2"

var (
	HTTPClientTimeout = 60
	client            = http.Client{Timeout: time.Duration(HTTPClientTimeout) * time.Second}
)

func GetPurchaseOrders(storage moyskladapptemplate.AppStorage, accountId string) (*PurchaseOrders, error) {
	accessToken, err := storage.AccessTokenByAccountId(accountId)
	if err != nil {
		return nil, err
	}

	req, err := utils.Request("GET", jsonEndpoint+"/entity/purchaseorder", accessToken, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var orders PurchaseOrders

	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&orders)
		if err != nil {
			return nil, err
		}
		return &orders, nil
	} else {
		var apiError moyskladapptemplate.JSONAPIError
		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			return nil, err
		}
		return nil, apiError
	}
}
