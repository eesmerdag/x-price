package router

import (
	priceclient "XPrice/web-app/price-client"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testPrices struct {
	Key      string
	Amount   int64
	Fraction int
}

func TestControllerGetPricesMethod(t *testing.T) {
	tests := []struct {
		name                      string
		req                       *http.Request
		successPricesFromServices []testPrices
		failPricesFromServices    []testPrices
		html                      string
	}{
		{
			name: "success from all services",
			req:  &http.Request{},
			successPricesFromServices: []testPrices{
				{Key: "first", Amount: 1923, Fraction: 2},
				{Key: "second", Amount: 9, Fraction: 1},
				{Key: "third", Amount: 1425, Fraction: 2},
			},
			html: "<h1>Macbook Air M2m Prices</h1>0.900000 <br>14.250000 <br>19.230000 <br>",
		},
		{
			name: "error from all services",
			req:  &http.Request{},
			failPricesFromServices: []testPrices{
				{Key: "first", Amount: 1923, Fraction: 2},
				{Key: "second", Amount: 9, Fraction: 1},
				{Key: "third", Amount: 1425, Fraction: 2},
			},
			html: "<h1>Macbook Air M2m Prices</h1>",
		},
		{
			name: "error from some services",
			req:  &http.Request{},
			successPricesFromServices: []testPrices{
				{Key: "first", Amount: 1923, Fraction: 2},
				{Key: "third", Amount: 1425, Fraction: 2},
			},
			failPricesFromServices: []testPrices{
				{Key: "second", Amount: 9, Fraction: 1},
			},
			html: "<h1>Macbook Air M2m Prices</h1>14.250000 <br>19.230000 <br>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var services []priceclient.PriceServiceClientI
			services = append(services, generateSuccessServices(tt.successPricesFromServices)...)
			services = append(services, generateFailServices(tt.failPricesFromServices)...)

			controller := NewController(services)

			res := http.HandlerFunc(controller.GetPrices)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)

			body, _ := ioutil.ReadAll(rr.Body)
			assert.True(t, strings.Compare(string(body), tt.html) == 0)
		})
	}
}

func generateSuccessServices(testPrices []testPrices) []priceclient.PriceServiceClientI {
	var priceServices []priceclient.PriceServiceClientI

	for _, p := range testPrices {
		var priceListClient priceclient.PriceServiceClientI
		priceListClient = *priceclient.NewMockPriceServiceClientSuccessCase(p.Key, p.Amount, p.Fraction)
		priceServices = append(priceServices, priceListClient)
	}

	return priceServices
}

func generateFailServices(testPrices []testPrices) []priceclient.PriceServiceClientI {
	var priceServices []priceclient.PriceServiceClientI

	for _, p := range testPrices {
		var priceListClient priceclient.PriceServiceClientI
		priceListClient = *priceclient.NewMockPriceServiceClientFailCase(p.Key)
		priceServices = append(priceServices, priceListClient)
	}

	return priceServices
}
