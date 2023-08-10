package main

import (
	priceclient "XPrice/web-app/price-client"
	"XPrice/web-app/router"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	c := &http.Client{
		Timeout: time.Second * 10,
	}
	var priceServices []priceclient.PriceServiceClientI

	// assume the first argument is for port and the rest for each price service
	if len(os.Args) > 2 {
		for i := 2; i <= len(os.Args)-1; i++ {
			var priceListClient priceclient.PriceServiceClientI
			priceListClient, err := priceclient.NewPriceServiceClient(c, os.Args[i])
			if err != nil {
				log.Fatalf("error some price services are not healty: " + os.Args[i])
			}
			priceServices = append(priceServices, priceListClient)
		}
	}

	web := router.NewController(priceServices)

	var port int64
	var err error
	if port, err = strconv.ParseInt(os.Args[1], 10, 64); err != nil {
		log.Fatalf("error initializing api: %v", err)
	}

	fmt.Println("Service is ready at port ", port)

	if err != http.ListenAndServe(":"+os.Args[1], web) {
		log.Fatalf("error initializing web app: %v", err)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
}
