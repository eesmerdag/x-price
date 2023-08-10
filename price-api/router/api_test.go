package router

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// as this is just a simple service and needed for code base, there is only one test case which is to retrieve price
func TestPriceEndpoint(t *testing.T) {
	test := struct {
		name  string
		req   *http.Request
		code  int
		price string
	}{
		name:  "ok",
		req:   &http.Request{},
		code:  http.StatusOK,
		price: "19.23",
	}

	t.Run(test.name, func(t *testing.T) {
		api, _ := NewApi(test.price)
		res := http.HandlerFunc(api.getMacbookPrice)
		rr := httptest.NewRecorder()
		res.ServeHTTP(rr, test.req)

		body, _ := ioutil.ReadAll(rr.Body)
		var priceInfo PriceInfo
		json.Unmarshal(body, &priceInfo)

		assert.True(t, rr.Code == test.code)
		assert.True(t, priceInfo.Display == "$19.23")
		assert.True(t, priceInfo.Currency == "USD")
		assert.True(t, priceInfo.Amount == 1923)
		assert.True(t, priceInfo.Fraction == 2)
	})
}
