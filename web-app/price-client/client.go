package price_client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// PriceServiceClientI this represents the clients from price-api, as all price-apis has the same format but different ports and response info, we can use just one client to handle all APIs
type PriceServiceClientI interface {
	GetPrice() (*PriceInfo, error)
	HealthCheck() error
	GetKey() string
}

type PriceServiceClient struct {
	httpClient HttpClient
	url        string
}

func NewPriceServiceClient(c *http.Client, url string) (*PriceServiceClient, error) {
	s := &PriceServiceClient{
		httpClient: c,
		url:        url,
	}
	var err error
	if err = s.HealthCheck(); err != nil {
		return nil, err
	}

	return s, nil
}

func (c PriceServiceClient) GetKey() string {
	return c.url
}

func (c PriceServiceClient) HealthCheck() error {
	req, err := http.NewRequest(http.MethodGet, c.url+"/status", nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("unhealthy service")
	}
	return nil
}

func (c PriceServiceClient) GetPrice() (*PriceInfo, error) {
	req, err := http.NewRequest(http.MethodGet, c.url+"/macbook-air-m2m/price", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, _ := ioutil.ReadAll(resp.Body)

	response := &PriceInfo{}

	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return response, nil
}
