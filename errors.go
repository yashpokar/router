package router

import "fmt"

type routeNotFoundError struct {
	path string
}

func (r *routeNotFoundError) Error() string {
	return fmt.Sprintf("route '%s' not found", r.path)
}
