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
	var variables []string

	for idx := range segments {
		segment := segments[idx]

		if root.children == nil {
			root.children = map[string]*tree{}
		}

		if isVariable(segment) {
			variable := strings.TrimLeft(segment, pathVariablePrefix)
			variables = append(variables, variable)
			segment = pathVariablePlaceholder
		}

		_, exists := root.children[segment]
		if !exists {
			root.children[segment] = &tree{}
		}
		root = root.children[segment]
	}

	route.variables = variables
	root.route = route
}

func (t *tree) find(path string) (*route, error) {
	root := t
	segments := strings.Split(path, pathSeparator)
	var values []string

	for idx := range segments {
		segment := segments[idx]

		_, exists := root.children[segment]
		if !exists {
			if _, exists = root.children[pathVariablePlaceholder]; !exists {
				return nil, &routeNotFoundError{path: path}
			}
			values = append(values, segment)
			segment = pathVariablePlaceholder
		}
		root = root.children[segment]
	}

	if root.route == nil {
		return nil, &routeNotFoundError{path: path}
	}

	root.route.variableValues = values
	return root.route, nil
}

func (t *tree) exists(path string) bool {
	root := t
	segments := strings.Split(path, pathSeparator)

	for idx := range segments {
		segment := segments[idx]

		_, exists := root.children[segment]
		if !exists {
			if _, exists = root.children[pathVariablePlaceholder]; !exists {
				return false
			}

			segment = pathVariablePlaceholder
		}
		root = root.children[segment]
	}

	return root.route != nil
}

func isVariable(s string) bool {
	return strings.HasPrefix(s, pathVariablePrefix)
}
