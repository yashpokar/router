package router

import "net/http"

type route struct {
	path    string
	handler http.HandlerFunc

	variables      []string
	variableValues []string
}

func newRoute(uri string, handler http.HandlerFunc) *route {
	return &route{
		path:    uri,
		handler: handler,
	}
}

func (r *route) getPathVariablesMap() map[string]string {
	if r.variables == nil {
		return nil
	}

	variables := map[string]string{}
	for idx := range r.variables {
		variables[r.variables[idx]] = r.variableValues[idx]
	}
	return variables
}
