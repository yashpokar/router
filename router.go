package router

import (
	"net/http"
	"strings"
)

type RouteRegistrar interface {
	GET(uri string, handler http.HandlerFunc)
	POST(uri string, handler http.HandlerFunc)
	PUT(uri string, handler http.HandlerFunc)
	PATCH(uri string, handler http.HandlerFunc)
	DELETE(uri string, handler http.HandlerFunc)
	Group(uri string) RouteRegistrar
}

type Router interface {
	RouteRegistrar
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}

type router struct {
	registry map[method]Tree
}

func createRegistry() map[method]Tree {
	return map[method]Tree{
		GET:    newTree(),
		POST:   newTree(),
		PUT:    newTree(),
		PATCH:  newTree(),
		DELETE: newTree(),
	}
}

func New() Router {
	return &router{
		registry: createRegistry(),
	}
}

func (r *router) GET(uri string, handler http.HandlerFunc) {
	r.register(GET, uri, handler)
}

func (r *router) POST(uri string, handler http.HandlerFunc) {
	r.register(POST, uri, handler)
}

func (r *router) PUT(uri string, handler http.HandlerFunc) {
	r.register(PUT, uri, handler)
}

func (r *router) PATCH(uri string, handler http.HandlerFunc) {
	r.register(PATCH, uri, handler)
}

func (r *router) DELETE(uri string, handler http.HandlerFunc) {
	r.register(DELETE, uri, handler)
}

func (r *router) Group(uri string) RouteRegistrar {
	return &group{basePath: uri, routerEngine: r}
}

func (r *router) register(m method, uri string, handler http.HandlerFunc) {
	uri = strings.Trim(uri, pathSeparator)

	r.registry[m].add(uri, newRoute(uri, handler))
}

func (r *router) resolve(m method, uri string) (*route, error) {
	return r.registry[m].find(uri)
}

func (r *router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	uri := strings.Trim(request.RequestURI, pathSeparator)
	m := stringToMethod(request.Method)

	route, err := r.resolve(m, uri)
	if err != nil {
		r.handleError(err, writer)
		return
	}

	route.handler(writer, request)
}

func (r *router) handleError(err error, writer http.ResponseWriter) {
	switch err.(type) {
	case *routeNotFoundError:
		writer.WriteHeader(http.StatusNotFound)
		return
	default:
		return
	}
}
