package request

import (
	"log"
	"net/http"
	"net/url"
)

// CreateRequest создает запрос с заранее подготовленными заголовками и добавляет в него переданные параметры (uriParams)
func CreateRequest(uri string, endpoint string, APIToken string, uriParams map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, uri+endpoint, nil)
	if err != nil {
		log.Fatalf("can't create random movie with Gosling request: %s", err)
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-API-KEY", APIToken)
	p := CreateParametresForRequest(uriParams)
	req.URL.RawQuery = p.Encode()
	return req, nil
}

// CreateParametresForRequest подготавливает параметры к передаче в запрос
func CreateParametresForRequest(p map[string]string) url.Values {
	params := url.Values{}
	for k, v := range p {
		params.Add(k, v)
	}
	return params
}
