package router

import (
	"net/http"
	"path"
	"strings"
)

type group struct {
	routerEngine *router
	basePath     string
}

func (g *group) register(m method, uri string, handler http.HandlerFunc) {
	uri = strings.Trim(path.Join(g.basePath, uri), pathSeparator)

	g.routerEngine.register(m, uri, handler)
}

func (g *group) GET(uri string, handler http.HandlerFunc) {
	g.register(GET, uri, handler)
}

func (g *group) POST(uri string, handler http.HandlerFunc) {
	g.register(POST, uri, handler)
}

func (g *group) PUT(uri string, handler http.HandlerFunc) {
	g.register(PUT, uri, handler)
}

func (g *group) PATCH(uri string, handler http.HandlerFunc) {
	g.register(PATCH, uri, handler)
}

func (g *group) DELETE(uri string, handler http.HandlerFunc) {
	g.register(DELETE, uri, handler)
}

func (g *group) Group(uri string) RouteRegistrar {
	uri = path.Join(g.basePath, uri)

	return &group{basePath: uri, routerEngine: g.routerEngine}
}
