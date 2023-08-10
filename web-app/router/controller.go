package router

import (
	priceclient "XPrice/web-app/price-client"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"sync"
	"time"
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

type PriceInMemory struct {
	info            *priceclient.PriceInfo
	lastUpdatedTime time.Time
}

type Controller struct {
	router        *mux.Router
	priceServices []priceclient.PriceServiceClientI
	mu            *sync.Mutex
	prices        map[string]*PriceInMemory // this is to keep prices in memory for some time. there is a cron job to delete related info from memory
}

func NewController(priceServices []priceclient.PriceServiceClientI) *Controller {
	router := mux.NewRouter()

	r := &Controller{
		router:        router,
		priceServices: priceServices,
		prices:        make(map[string]*PriceInMemory),
		mu:            &sync.Mutex{},
	}

	router.HandleFunc("/prices", r.GetPrices).Methods(http.MethodGet)
	router.Use(panicRecovery)
	return r
}

func (a Controller) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.router.ServeHTTP(w, req)
}

func (a Controller) addRespToMap(key string, resp *priceclient.PriceInfo) {
	a.mu.Lock()
	a.prices[key] = &PriceInMemory{info: resp, lastUpdatedTime: time.Now()}
	a.mu.Unlock()
}

func (a Controller) getRespFromMap(key string) (resp *PriceInMemory, ok bool) {
	a.mu.Lock()
	resp, ok = a.prices[key]
	a.mu.Unlock()
	return resp, ok
}

func (a Controller) GetPrices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	ch := make(chan *priceclient.PriceInfo, len(a.priceServices))
	var wg sync.WaitGroup
	for i := 0; i < len(a.priceServices); i++ {
		wg.Add(1)
		go a.priceInfoFromService(a.priceServices[i], ch, &wg)
	}

	wg.Wait()

	var prices []*priceclient.PriceInfo
	for i := 0; i < len(a.priceServices); i++ {
		result := <-ch
		// amount -1 means there was error with one service
		if result.Amount != -1 {
			prices = append(prices, result)
		}
	}

	fmt.Fprintf(w, "<h1>Macbook Air M2m Prices</h1>"+sortPrices(prices))
}

func (a Controller) priceInfoFromService(client priceclient.PriceServiceClientI, ch chan *priceclient.PriceInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	// at this point, the service should have caching mechanism since the price is not updated frequently at all.
	// however, i just assume there is no caching mechanism in APIs itself so as a solution, i decided to keep the info into memory of Web Application
	// so that we don't need to call request to other services
	val, ok := a.getRespFromMap(client.GetKey())
	// update value if previous one 1 min ago
	if ok && val.lastUpdatedTime.After(time.Now().Add(time.Minute*-1)) {
		ch <- val.info
		return
	}

	resp, err := client.GetPrice()
	if err != nil {
		ch <- &priceclient.PriceInfo{
			Amount: -1, // means some error happened at this service so we won't show this in the list
		}
		return
	}
	ch <- resp

	a.addRespToMap(client.GetKey(), resp)
}

func sortPrices(prices []*priceclient.PriceInfo) string {
	var amounts []float64
	for i := 0; i < len(prices); i++ {
		amounts = append(amounts, float64(prices[i].Amount)/float64(prices[i].Fraction))
	}
	sort.Float64s(amounts)

	var message = ""
	for i := 0; i < len(amounts); i++ {
		message = message + "<div>" + fmt.Sprintf("%f", amounts[i]) + "<div> <br>"
	}
	return message
}

func panicRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "internal server error")
				return
			}
		}()

		h.ServeHTTP(w, r)
	})
}
