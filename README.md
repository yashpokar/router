# Router [![Build](https://github.com/yashpokar/router/actions/workflows/ci.yaml/badge.svg?branch=master)](https://github.com/yashpokar/router/actions/workflows/ci.yaml)
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
