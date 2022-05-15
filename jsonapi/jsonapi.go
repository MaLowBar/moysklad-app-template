package jsonapi

import (
	"encoding/json"
	"fmt"
	templ "github.com/MaLowBar/moysklad-app-template"
	"github.com/MaLowBar/moysklad-app-template/utils"
	"net/http"
	"time"
)

/*
	TODO: сделать получение информации о
 		- заказах поставщикам (complete),
 		- счетах поставщиков,
 		- приемках,
 		- заказах покупателей,
		- счетах покупателям,
 		- отгрузках
*/

const jsonEndpoint = "https://online.moysklad.ru/api/remap/1.2"

var (
	HTTPClientTimeout = 60
	client            = http.Client{Timeout: time.Duration(HTTPClientTimeout) * time.Second}
)

func getEntity[T any](url, accessToken string) (*T, error) {
	req, err := utils.Request("GET", url, accessToken, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		ent := new(T)
		ent, err = utils.Fetch[*T](resp.Body)
		if err != nil {
			return nil, err
		}
		return ent, nil
	} else {
		apiErr, err := utils.Fetch[templ.JSONAPIError](resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, apiErr
	}
}

type metaWrap struct {
	Meta Meta `json:"meta"`
}

func GetAllEntities[T any](storage templ.AppStorage, accountId, entity string) ([]T, error) {
	accessToken, err := storage.AccessTokenByAccountId(accountId)
	if err != nil {
		return nil, err
	}

	url := jsonEndpoint + "/entity/" + entity

	res := make([]T, 0)
	ent, err := getEntity[T](url, accessToken)
	if err != nil {
		return nil, err
	}
	res = append(res, *ent)

	meta, err := getEntity[metaWrap](url, accessToken)
	if err != nil {
		return nil, err
	}

	pages := meta.Meta.Size / 1000
	if meta.Meta.Size <= 1000 {
		return res, nil
	}

	input := make(chan int, pages)
	errors := make(chan error, pages)
	for i := 1; i <= pages; i++ {
		input <- i * 1000
	}
	close(input)

	maxWorkers := 5
	for i := 0; i < maxWorkers; i++ {
		go func(input chan int, errors chan error) {
			for offset := range input {
				ent, err := getEntity[T](fmt.Sprintf("%s?offset=%d", url, offset), accessToken)
				if err != nil {
					errors <- err
					continue
				}
				res = append(res, *ent)
				errors <- nil
			}
		}(input, errors)
	}

	for i := 0; i < len(errors); i++ {
		if err = <-errors; err != nil {
			return nil, err
		}
	}
	return res, nil
}

func GetPurchaseOrders(storage templ.AppStorage, accountId string) (*PurchaseOrders, error) {
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
		var apiError templ.JSONAPIError
		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			return nil, err
		}
		return nil, apiError
	}
}
