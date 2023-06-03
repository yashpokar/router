package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yashpokar/router/tests"
)

func TestTree(t *testing.T) {
	indexRoute := &route{
		path:    "/",
		handler: tests.Handler,
	}

	t.Run("creates tree", func(t *testing.T) {
		tree := newTree()
		tree.add("/", indexRoute)

		r, err := tree.find("/")

		assert.NoError(t, err)
		assert.Equal(t, r, indexRoute)
	})

	t.Run("returns error when no route define", func(t *testing.T) {
		tree := newTree()
		tree.add("/", indexRoute)

		r, err := tree.find("/non-existing-path")

		assert.EqualError(t, err, "route '/non-existing-path' not found")
		assert.Nil(t, r)
	})

	t.Run("returns true when route is define", func(t *testing.T) {
		tree := newTree()
		tree.add("/index", indexRoute)

		exists := tree.exists("/index")
		assert.True(t, exists)
	})

	t.Run("returns true when route is define (nested path)", func(t *testing.T) {
		tree := newTree()
		tree.add("/nested/path", indexRoute)

		exists := tree.exists("/nested/path")
		assert.True(t, exists)
	})

	t.Run("returns route not found error when route is not define", func(t *testing.T) {
		tree := newTree()

		route, err := tree.find("/nested/path")
		assert.EqualError(t, err, "route '/nested/path' not found")
		assert.Nil(t, route)
	})

	t.Run("returns route (path variable)", func(t *testing.T) {
		tree := newTree()
		tree.add("/orders/:order_id", indexRoute)

		route, err := tree.find("/orders/147db1c8-7ff2-4740-a6ba-0ba5761224ca")
		assert.NoError(t, err)
		assert.Equal(t, route, indexRoute)
	})

	t.Run("returns false when route is not define", func(t *testing.T) {
		tree := newTree()

		exists := tree.exists("/nested/path")
		assert.False(t, exists)
	})

	t.Run("returns true when route is define (path variable)", func(t *testing.T) {
		tree := newTree()
		tree.add("/api/v1/users/:user_id", &route{
			path:    "/api/v1/users/:user_id",
			handler: tests.Handler,
		})

		exists := tree.exists("/api/v1/users/jhon")
		assert.True(t, exists)
	})

	t.Run("returns false when variable path does not match", func(t *testing.T) {
		tree := newTree()
		tree.add("/api/v1/users/:user_id/edit", &route{
			path:    "/api/v1/users/:user_id/edit",
			handler: tests.Handler,
		})

		exists := tree.exists("/api/v1/users/jhon")
		assert.False(t, exists)
	})

	t.Run("returns false when tree has nodes but the require node is not the end", func(t *testing.T) {
		tree := newTree()
		tree.add("/another/level/nested/path", indexRoute)

		route, err := tree.find("/another/level")
		assert.EqualError(t, err, "route '/another/level' not found")
		assert.Nil(t, route)
	})
}
