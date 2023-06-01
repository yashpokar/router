package router

import "net/http"

type route struct {
	path    string
	handler http.HandlerFunc
}

func newRoute(uri string, handler http.HandlerFunc) *route {
	return &route{
		path:    uri,
		handler: handler,
	}
}
