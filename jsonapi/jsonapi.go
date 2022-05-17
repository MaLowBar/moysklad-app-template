package jsonapi

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	templ "github.com/MaLowBar/moysklad-app-template"
	"github.com/MaLowBar/moysklad-app-template/utils"
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
	Tryes             = 3
	Timeout           = 3
)

type Entities[T any] struct {
	Meta    Meta `json:"meta"`
	Context struct {
		Employee struct {
			Meta Meta `json:"meta"`
		} `json:"employee"`
	} `json:"context"`
	Rows []T `json:"rows"`
}

func getEntity[T any](url, accessToken string) (*T, error) {
	res := new(T)
	var err error
	tryCount := Tryes
	for tryCount > 0 {
		res, err = getEnt[T](url, accessToken)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(Timeout) * time.Second)
		tryCount--
	}
	return res, err
}

func getEnt[T any](url, accessToken string) (*T, error) {
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

func GetAllEntities[T any, _ Entities[T]](storage templ.AppStorage, accountId, entity string) ([]T, error) {
	accessToken, err := storage.AccessTokenByAccountId(accountId)
	if err != nil {
		return nil, err
	}

	url := jsonEndpoint + "/entity/" + entity

	res := make([]T, 0)

	meta, err := getEntity[metaWrap](fmt.Sprintf("%s?limit=1", url), accessToken)
	if err != nil {
		return nil, err
	}

	pages := 1 + meta.Meta.Size/1000

	input := make(chan int, pages)
	errors := make(chan error, pages)
	for i := 0; i < pages; i++ {
		input <- i * 1000
	}
	close(input)
	wg := sync.WaitGroup{}
	maxWorkers := 5
	wg.Add(maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go func(input chan int, errors chan error) {
			for offset := range input {
				ent, err := getEntity[Entities[T]](fmt.Sprintf("%s?offset=%d", url, offset), accessToken)
				if err != nil {
					errors <- err
					continue
				}
				res = append(res, ent.Rows...)
				errors <- nil
			}
			wg.Done()
		}(input, errors)
	}
	wg.Wait()
	close(errors)

	for err = range errors {
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
