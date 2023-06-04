# Router [![Build](https://github.com/yashpokar/router/actions/workflows/ci.yaml/badge.svg)](https://github.com/yashpokar/router/actions/workflows/ci.yaml)
Go based lightweight http router.

## Basic Router

```go
package app

import (
    "net/http"

    "github.com/yashpokar/router"
)

func BasicHandler(w http.ResponseWriter, r *http.Request) {
    // TODO : do cool stuff
}

func main() {
    r := router.New()
	r.GET("/", BasicHandler)

    server := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }
    server.ListenAndServe()
}
```

## Grouping

```go
r := router.New()
r.GET("/", Index)

api := r.Group("/api")
v1 := api.Group("/v1")

orders := v1.Group("/orders")
orders.GET("/", ListOrders)
orders.POST("/", PlaceOrder)
orders.GET("/:orderId", GetOrder)
```

## Path Variables

```go
func GetOrder(w http.ResponseWriter, r *http.Request) {
    vars := router.Vars(r)
    orderID := vars["orderId"]

    // TODO : fetch order details
    // TODO : send response
}

r := router.New()

r.GET("/orders/:orderId", GetOrder)
```

## Middleware

The router does not provide utility for middlewares because implementing raw middlewares are very easy
and unnecessary abstraction can slow down your application.

Here's an example how you can define your own middlewares.

```go

func IPCollectorMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request) {
        ip := r.Header.Get("X-Forwarded-For")
        fmt.Printf("The IP address of source is %s", ip)

        next.ServeHTTP(writer, request)
    }
}

func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request) {
		authorizationHeader := request.Header.Get("Authorization")
		if authorizationHeader == "" {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        token := strings.TrimRight(authorizationHeader, "Bearer ")

        // ...

        next.ServeHTTP(writer, request)
    }
}

func PrivateHandler(w http.ResponseWriter, r *http.Request) {
    // ...
}

r := router.New()

r.GET("/private", IPCollectorMiddleware(AuthenticationMiddleware(PrivateHandler)))
```
