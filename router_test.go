package router

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yashpokar/router/tests"
)

func TestRouter(t *testing.T) {
	r := New().(*router)
	server := http.Server{Handler: r, Addr: ":8787"}
	go func() {
		err := server.ListenAndServe()
		assert.NoError(t, err)
	}()

	t.Run("registers GET method", func(t *testing.T) {
		r.GET("/", tests.Handler)
		route, err := r.resolve("GET", "")
		assert.NoError(t, err)
		assert.Equal(t, route.path, "")
	})

	t.Run("registers POST method", func(t *testing.T) {
		r.POST("/", tests.Handler)
		route, err := r.resolve("POST", "")
		assert.NoError(t, err)
		assert.Equal(t, route.path, "")
	})

	t.Run("registers PUT method", func(t *testing.T) {
		r.PUT("/", tests.Handler)
		route, err := r.resolve("PUT", "")
		assert.NoError(t, err)
		assert.Equal(t, route.path, "")
	})

	t.Run("registers PATCH method", func(t *testing.T) {
		r.PATCH("/", tests.Handler)
		route, err := r.resolve("PATCH", "")
		assert.NoError(t, err)
		assert.Equal(t, route.path, "")
	})

	t.Run("registers DELETE method", func(t *testing.T) {
		r.DELETE("/", tests.Handler)
		route, err := r.resolve("DELETE", "")
		assert.NoError(t, err)
		assert.Equal(t, route.path, "")
	})

	t.Run("handles general error", func(t *testing.T) {
		r.handleError(errors.New("internal server error"), tests.NewMockResponseWriter())
	})

	t.Run("resolves the route", func(t *testing.T) {
		r.GET("/path/nested", tests.Handler)
		route, err := r.resolve("GET", "path/nested")

		assert.NoError(t, err)
		assert.Equal(t, "path/nested", route.path)
	})
}
