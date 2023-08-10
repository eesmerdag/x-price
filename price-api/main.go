package main

import (
	"XPrice/price-api/router"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	// assuming service run is called with two arguments where first one represents port and second one is price of product
	amount := os.Args[2]
	api, err := router.NewApi(amount)
	if err != nil {
		log.Fatalf("error setting api up: %v", err)
	}

	var port int64
	if port, err = strconv.ParseInt(os.Args[1], 10, 64); err != nil {
		log.Fatalf("error initializing api: %v", err)
	}

	fmt.Println("Service is ready at port ", port)
	if err != http.ListenAndServe(":"+os.Args[1], api) {
		log.Fatalf("error initializing api: %v", err)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
}
