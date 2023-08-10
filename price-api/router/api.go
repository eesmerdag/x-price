package router

import (
	"encoding/json"
	"github.com/Rhymond/go-money"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ErrorResp struct {
	Message string
	Code    int
}

type PriceInfo struct {
	Display  string
	Currency string
	Amount   int64
	Fraction int
}

type Api struct {
	router *mux.Router
	price  float64
}

func NewApi(price string) (*Api, error) {
	var pr float64
	var err error
	if pr, err = strconv.ParseFloat(price, 64); err != nil {
		return nil, err
	}

	router := mux.NewRouter()

	r := &Api{
		router: router,
		price:  pr,
	}

	router.HandleFunc("/macbook-air-m2m/price", r.getMacbookPrice).Methods(http.MethodGet)
	router.HandleFunc("/status", r.status).Methods(http.MethodGet)
	router.Use(panicRecovery)
	return r, nil
}

func (a Api) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.router.ServeHTTP(w, req)
}

func (a Api) getMacbookPrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	price := money.NewFromFloat(a.price, money.USD)

	res := PriceInfo{
		Display:  price.Display(),
		Currency: price.Currency().Code,
		Amount:   price.Amount(),
		Fraction: price.Currency().Fraction,
	}
	json.NewEncoder(w).Encode(res)
}

func (a Api) status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func panicRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errorResponse(w, "unexpected internal error", http.StatusInternalServerError)
				return
			}
		}()

		h.ServeHTTP(w, r)
	})
}

func errorResponse(w http.ResponseWriter, message string, code int) {
	errObj := ErrorResp{
		Message: message,
		Code:    code,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errObj)
}
