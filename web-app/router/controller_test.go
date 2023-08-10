package router

import (
	priceclient "XPrice/web-app/price-client"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestControllerGetPricesMethod(t *testing.T) {
	test := struct {
		name          string
		req           *http.Request
		code          int
		priceServices []priceclient.PriceServiceClientI
	}{
		name:          "ok",
		req:           &http.Request{},
		code:          http.StatusOK,
		priceServices: generateServices(),
	}

	t.Run(test.name, func(t *testing.T) {
		/*		mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()

				mockPriceServiceClient := priceclient.NewMockPriceServiceClientI(mockCtrl)

				mockPriceServiceClient.EXPECT().GetPrice().Return()

				controller := NewController(test.priceServices)
		*/
		controller := NewController(test.priceServices)

		res := http.HandlerFunc(controller.GetPrices)
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

func generateServices() []priceclient.PriceServiceClientI {
	c := &http.Client{
		Timeout: time.Second * 10,
	}
	var priceServices []priceclient.PriceServiceClientI

	var priceListClient priceclient.PriceServiceClientI
	priceListClient, _ = priceclient.NewPriceServiceClient(c, "test:1001")
	priceServices = append(priceServices, priceListClient)

	return priceServices
}
