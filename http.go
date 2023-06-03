package router

import "fmt"

type method string

const (
	GET    method = "GET"
	POST   method = "POST"
	PUT    method = "PUT"
	PATCH  method = "PATCH"
	DELETE method = "DELETE"
)

func stringToMethod(name string) method {
	switch name {
	case "GET":
		return GET
	case "POST":
		return POST
	case "PUT":
		return PUT
	case "PATCH":
		return PATCH
	case "DELETE":
		return DELETE
	default:
		panic(fmt.Sprintf("method %s is not supported", name))
	}
}
