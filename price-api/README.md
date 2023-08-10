# price-api

This API is to provide price info for macbook air m2m as a data source.

# how to run the service

Service needs two different arguments while running. These are port and price info. In this way, we would be able to
have many different sources with same code piece.

To run the service, go to project folder and hit:

go run price-api/main.go PORT PRICE

example:

```
go run price-api/main.go  1923 19.23
```

to run the test:

```
go test ./price-api/router/
```

Assumption is that prices are based on USD. example response:

```
{
    "Display": "$19.23",
    "Currency": "USD",
    "Amount": 1923,
    "Fraction": 2
}
```