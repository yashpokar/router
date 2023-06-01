package router

import (
	"strings"
)

type Tree interface {
	add(path string, route *route)
	find(path string) (*route, error)
	exists(path string) bool
}

type tree struct {
	children map[string]*tree
	route    *route
}

func newTree() Tree {
	return &tree{}
}

func (t *tree) add(path string, route *route) {
	root := t
	segments := strings.Split(path, pathSeparator)

	for idx := range segments {
		segment := segments[idx]

		if root.children == nil {
			root.children = map[string]*tree{}
		}

		_, exists := root.children[segment]
		if !exists {
			root.children[segment] = &tree{}
		}

		root = root.children[segment]
	}

	root.route = route
}

func (t *tree) find(path string) (*route, error) {
	root := t
	segments := strings.Split(path, pathSeparator)

	for idx := range segments {
		segment := segments[idx]

		_, exists := root.children[segment]
		if !exists {
			return nil, &routeNotFoundError{path: path}
		}

		root = root.children[segment]
	}

	if root.route == nil {
		return nil, &routeNotFoundError{path: path}
	}

	return root.route, nil
}

func (t *tree) exists(path string) bool {
	root := t
	segments := strings.Split(path, pathSeparator)

	for idx := range segments {
		segment := segments[idx]

		_, exists := root.children[segment]
		if !exists {
			return false
		}

		root = root.children[segment]
	}

	return root.route != nil
}
