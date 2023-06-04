package router

import (
	"context"
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

type PanicHandler func(w http.ResponseWriter, r *http.Request, err any)

type Router interface {
	RouteRegistrar
	OnPanic(handler PanicHandler)
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}

type router struct {
	registry     map[method]Tree
	panicHandler PanicHandler
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
	defer func() {
		if e := recover(); e != nil && r.panicHandler != nil {
			r.panicHandler(writer, request, e)
		}
	}()

	uri := strings.Trim(request.RequestURI, pathSeparator)
	m := stringToMethod(request.Method)

	route, err := r.resolve(m, uri)
	if err != nil {
		r.handleError(err, writer)
		return
	}

	pathVariables := route.getPathVariablesMap()
	ctx := request.Context()
	ctx = context.WithValue(ctx, pathVariablesContextKey, pathVariables)
	request = request.WithContext(ctx)
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

func (r *router) OnPanic(handler PanicHandler) {
	r.panicHandler = handler
}
