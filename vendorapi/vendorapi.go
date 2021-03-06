package vendorapi

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	moyskladapptemplate "github.com/MaLowBar/moysklad-app-template"
	"github.com/MaLowBar/moysklad-app-template/utils"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
)

const marketplaceEndpoint = "https://online.moysklad.ru/api/vendor/1.0"

var (
	HTTPClientTimeout = uint64(60)
	client            = http.Client{Timeout: time.Duration(HTTPClientTimeout) * time.Second}
)

func vendorRequest(appInfo moyskladapptemplate.AppConfig, method, url string, body io.Reader) (*http.Request, error) {
	someBytes := make([]byte, 16)
	if _, err := rand.Read(someBytes); err != nil {
		return nil, err
	}

	jti, err := uuid.FromBytes(someBytes)
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:  appInfo.UID,
		IssuedAt: time.Now().UTC().Unix(),
		Id:       jti.String(),
	})

	tokenString, err := token.SignedString([]byte(appInfo.SecretKey))
	if err != nil {
		return nil, err
	}

	return utils.Request(method, url, tokenString, body)
}

func GetUserContext(contextKey string, appInfo moyskladapptemplate.AppConfig) (*UserContext, error) {
	req, err := vendorRequest(appInfo, "POST", marketplaceEndpoint+"/context/"+contextKey, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userContext UserContext

	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&userContext)
		if err != nil {
			return nil, err
		}
		return &userContext, nil
	} else {
		var body []byte
		_, err = resp.Body.Read(body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("wrong status code %d while getting user context, repsponse body: %s", resp.StatusCode, body)

	}

}
