# price-api

This is simple Web application that shows prices of Macbook Air M2M

# how to run the service

Service needs port arguments as the first one and many other optional services endpoints. There is not limit for services which provides macbook price. But the all of the services should be run before this particular one.

To run the service, go to project folder and hit:

go run price-api/main.go PORT SERVICE1 SERVICE2 ...

example:

```
go run web-app/main.go 1933 http://127.0.0.1:1903 http://127.0.0.1:1904  http://127.0.0.1:1905 http://127.0.0.1:1902 http://127.0.0.1:1901 http://127.0.0.1:1900 http://127.0.0.1:1899
```

to run the test:

```
go test ./web-app/router/
```

Sample web page for prices page at /prices 

![Screen Shot 2023-08-10 at 13.06.32.png](..%2F..%2F..%2FDesktop%2FScreen%20Shot%202023-08-10%20at%2013.06.32.png)